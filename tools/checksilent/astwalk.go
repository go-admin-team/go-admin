package main

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// structLiteral is a composite literal whose type resolved to a named struct.
type structLiteral struct {
	PkgPath string
	Name    string
	Lit     *ast.CompositeLit
}

// forEachStructLiteral visits every composite literal in the file whose type
// resolves to a named type, including the ones written with the type elided.
//
// The elided form is the one that matters: seed data is written as
// []models.SysMenu{{MenuId: 9000}, {MenuId: 9001}}, and the inner literals carry
// no type of their own. A walker that only looked at CompositeLit.Type would
// silently skip every seed in the repository and report nothing, which for a
// tool about silent failure would be its own punchline.
func forEachStructLiteral(sf *sourceFile, fn func(structLiteral)) {
	// The type each type-less literal inherits from the literal containing it.
	elided := map[*ast.CompositeLit]ast.Expr{}

	var propagate func(lit *ast.CompositeLit, typ ast.Expr)
	propagate = func(lit *ast.CompositeLit, typ ast.Expr) {
		child := elementType(typ)
		if child == nil {
			return
		}
		for _, elt := range lit.Elts {
			v := elt
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				v = kv.Value
			}
			if cl, ok := v.(*ast.CompositeLit); ok && cl.Type == nil {
				elided[cl] = child
				propagate(cl, child)
			}
		}
	}

	// ast.Inspect visits a node before its children, so every typed literal
	// fills in its descendants before the reporting pass reaches them.
	ast.Inspect(sf.Syntax, func(n ast.Node) bool {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		if typ := litType(cl, elided); typ != nil {
			propagate(cl, typ)
		}
		return true
	})

	ast.Inspect(sf.Syntax, func(n ast.Node) bool {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		typ := litType(cl, elided)
		if typ == nil {
			return true
		}
		if pkg, name, ok := resolveNamed(sf, typ); ok {
			fn(structLiteral{PkgPath: pkg, Name: name, Lit: cl})
		}
		return true
	})
}

func litType(cl *ast.CompositeLit, elided map[*ast.CompositeLit]ast.Expr) ast.Expr {
	if cl.Type != nil {
		return cl.Type
	}
	return elided[cl]
}

// resolveNamed maps a type expression to (import path, type name). A bare
// identifier means a type declared in this file's own package.
func resolveNamed(sf *sourceFile, typ ast.Expr) (string, string, bool) {
	switch t := typ.(type) {
	case *ast.StarExpr:
		return resolveNamed(sf, t.X)
	case *ast.Ident:
		return sf.Pkg, t.Name, true
	case *ast.SelectorExpr:
		pkgIdent, ok := t.X.(*ast.Ident)
		if !ok {
			return "", "", false
		}
		path, ok := sf.imports[pkgIdent.Name]
		if !ok {
			return "", "", false
		}
		return path, t.Sel.Name, true
	}
	return "", "", false
}

// elementType is the type the children of a composite literal take when they
// leave theirs out.
func elementType(typ ast.Expr) ast.Expr {
	switch t := typ.(type) {
	case *ast.ArrayType:
		return t.Elt
	case *ast.MapType:
		return t.Value
	case *ast.StarExpr:
		return elementType(t.X)
	}
	return nil
}

// field returns the value written for a named field of a struct literal.
func field(lit *ast.CompositeLit, name string) (ast.Expr, bool) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if key, ok := kv.Key.(*ast.Ident); ok && key.Name == name {
			return kv.Value, true
		}
	}
	return nil, false
}

// intValue evaluates an integer field: a literal, a negated literal, or an
// identifier naming a constant in the same package.
//
// Anything computed at run time is skipped rather than guessed at. That is the
// one thing these checks miss, and missing is the right way to be wrong here -
// a false positive teaches people to add ignore comments, and then the tool is
// finished.
func intValue(sf *sourceFile, expr ast.Expr) (int64, bool) {
	switch e := expr.(type) {
	case *ast.Ident:
		v, ok := sf.consts[e.Name]
		return v, ok
	case *ast.UnaryExpr:
		if e.Op == token.SUB {
			if v, ok := intValue(sf, e.X); ok {
				return -v, true
			}
		}
	}
	return intLiteral(expr)
}

func intLiteral(expr ast.Expr) (int64, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, false
	}
	v, err := strconv.ParseInt(strings.ReplaceAll(lit.Value, "_", ""), 0, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// stringValue evaluates a string field: a literal, a concatenation of literals,
// or an identifier naming a string constant in the same package.
func stringValue(sf *sourceFile, expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(e.Value)
		if err != nil {
			return "", false
		}
		return s, true
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}
		l, lok := stringValue(sf, e.X)
		r, rok := stringValue(sf, e.Y)
		if lok && rok {
			return l + r, true
		}
	}
	return "", false
}

// embeddedTypes returns the types a struct embeds, as (import path, name).
func embeddedTypes(sf *sourceFile, st *ast.StructType) [][2]string {
	var out [][2]string
	for _, f := range st.Fields.List {
		if len(f.Names) != 0 {
			continue // a named field, not an embed
		}
		if pkg, name, ok := resolveNamed(sf, f.Type); ok {
			out = append(out, [2]string{pkg, name})
		}
	}
	return out
}

// tableNames maps struct name to the literal its TableName method returns.
func tableNames(sf *sourceFile) map[string]string {
	out := map[string]string{}
	for _, decl := range sf.Syntax.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "TableName" || fn.Recv == nil || len(fn.Recv.List) != 1 || fn.Body == nil {
			continue
		}
		recv := receiverName(fn.Recv.List[0].Type)
		if recv == "" {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			ret, ok := n.(*ast.ReturnStmt)
			if !ok || len(ret.Results) != 1 {
				return true
			}
			if s, ok := stringValue(sf, ret.Results[0]); ok && out[recv] == "" {
				out[recv] = s
			}
			return true
		})
	}
	return out
}

func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverName(t.X)
	}
	return ""
}

// structTypes maps struct name to its declaration.
func structTypes(sf *sourceFile) map[string]*ast.StructType {
	out := map[string]*ast.StructType{}
	ast.Inspect(sf.Syntax, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if st, ok := ts.Type.(*ast.StructType); ok {
			out[ts.Name.Name] = st
		}
		return true
	})
	return out
}
