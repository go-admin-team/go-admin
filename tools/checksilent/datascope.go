package main

import (
	"go/ast"
	"strings"
)

// ---------------------------------------------------------------------------
// check 8: a handler that reads the data permission, on a route that never
// installs the middleware which puts one there
// ---------------------------------------------------------------------------

// permissionGetter is the function a handler calls to obtain the caller's data
// scope, and permissionMiddleware is the middleware that puts one in the
// context. Matched by name rather than by resolved symbol: the tool parses
// without type checking, and both names are distinctive enough that a
// same-named function from somewhere else would still be worth a look.
const (
	permissionGetter     = "GetPermissionFromContext"
	permissionMiddleware = "PermissionAction"
)

// actionsPkgSuffix identifies the package the two names above live in - this
// repository's common/actions shim and core's sdk/contract/actions both end
// this way, and a module rename changes neither.
const actionsPkgSuffix = "/actions"

// handlerKey identifies one handler method uniquely across packages, so that
// two types named SysUser in different packages are not confused.
type handlerKey struct {
	Pkg  string
	Type string
	Func string
}

// checkDataScopeRoutes reports a route whose handler asks for the caller's data
// permission while the group it is registered on never installs the middleware
// that supplies one.
//
// GetPermissionFromContext cannot fail. When nothing put a *DataPermission in
// the context it hands back a zero value, whose DataScope is the empty string -
// and the empty string is not one of the five scopes Permission recognises, so
// it takes the default branch. That branch fails closed: the query is given
// `1 = 0` and matches nothing.
//
// The result is an endpoint that answers "not found" or "no permission" for
// rows that plainly exist, and only on deployments that set enabledp: true -
// with data permissions off, Permission returns the query untouched and the
// missing middleware costs nothing. That is the shape this check exists for: a
// default configuration where the mistake is invisible, and a test suite that
// runs on it.
//
// It happened. /api/v1/getinfo read the permission on a group carrying only the
// JWT middleware, so every login on a deployment with data permissions enabled
// ended in a 401 from the endpoint the browser calls immediately after signing
// in - and went back to the login page.
//
// Either half is a fix, and which one depends on the route. A handler that
// reads somebody else's rows wants the middleware. A handler reading the
// caller's own row - where the id comes from the token - wants no scope at all,
// because a scope has nothing left to restrict there and DataScopeSelf, which
// matches on create_by, would reject every user who did not create their own
// account. The check reports the mismatch and leaves the choice.
func checkDataScopeRoutes(s *snapshot) []Finding {
	handlers := permissionReadingHandlers(s)
	if len(handlers) == 0 {
		return nil
	}

	var out []Finding
	for _, sf := range s.Files {
		if sf.isTest() {
			continue
		}
		for _, decl := range sf.Syntax.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			out = append(out, s.routeFindings(sf, fn, handlers)...)
		}
	}
	return out
}

// permissionReadingHandlers collects every method whose body calls the getter.
//
// Test files are included deliberately: a handler is a handler wherever it is
// declared, and skipping them would let a route registered from a test fixture
// go unchecked while the fixture is exactly where a new one gets written first.
func permissionReadingHandlers(s *snapshot) map[handlerKey]bool {
	out := map[handlerKey]bool{}
	for _, sf := range s.Files {
		for _, decl := range sf.Syntax.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			recv := receiverTypeName(fn.Recv.List[0].Type)
			if recv == "" {
				continue
			}
			if callsPackageFunc(sf, fn.Body, permissionGetter) {
				out[handlerKey{Pkg: sf.Pkg, Type: recv, Func: fn.Name.Name}] = true
			}
		}
	}
	return out
}

// routeFindings walks one function looking for group definitions and the routes
// registered on them.
func (s *snapshot) routeFindings(sf *sourceFile, fn *ast.FuncDecl, handlers map[handlerKey]bool) []Finding {
	// Local variable bindings, filled as the body is walked in order so that a
	// registration only ever sees definitions that precede it.
	apiVars := map[string]handlerKey{} // var -> the type it holds
	guarded := map[string]bool{}       // group var -> middleware installed
	known := map[string]bool{}         // group var -> is a router group at all
	prefix := map[string]string{}      // group var -> the path it was declared with

	var out []Finding
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			for i, lhs := range stmt.Lhs {
				id, ok := lhs.(*ast.Ident)
				if !ok || i >= len(stmt.Rhs) {
					continue
				}
				rhs := stmt.Rhs[i]
				if key, ok := apiTypeOf(sf, rhs); ok {
					apiVars[id.Name] = key
					continue
				}
				if parent, isGroup := groupSource(rhs); isGroup {
					known[id.Name] = true
					prefix[id.Name] = prefix[parent] + groupPath(rhs)
					// A subgroup inherits whatever its parent already had:
					// gin copies the parent's handler chain into the child.
					guarded[id.Name] = guarded[parent] || containsCallNamed(rhs, permissionMiddleware)
				}
			}
		case *ast.ExprStmt:
			// A separate `g.Use(...)` after the group was defined.
			call, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				return true
			}
			if target, ok := receiverIdentOf(call, "Use"); ok && known[target] {
				if containsCallNamed(call, permissionMiddleware) {
					guarded[target] = true
				}
			}
		}
		return true
	})

	// Second pass for the registrations, so that a `.Use` written below a route
	// still counts - the middleware chain is assembled before any request is
	// served, not in source order.
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		gvar, method, ok := routeRegistration(call)
		if !ok || !known[gvar] || guarded[gvar] {
			return true
		}
		route, handlerVar, handlerName, ok := routeArgs(call)
		if !ok {
			return true
		}
		key, ok := apiVars[handlerVar]
		if !ok {
			return true
		}
		key.Func = handlerName
		if !handlers[key] {
			return true
		}
		out = append(out, s.finding(Error, checkDataScopeRoute, sf, call,
			"%s %q is handled by %s.%s, which reads the caller's data permission,\n"+
				"    but the group it is registered on never installs %s.\n"+
				"    GetPermissionFromContext then returns the zero value, whose empty DataScope is not a\n"+
				"    recognised scope, so Permission fails closed and the query matches nothing - on any\n"+
				"    deployment with enabledp: true. With data permissions off the route works, which is\n"+
				"    why this does not show up in the default configuration or in CI.\n"+
				"    Add %s() to the group, or stop scoping a query that is already limited to the caller.",
			method, prefix[gvar]+route, key.Type, handlerName, permissionMiddleware, permissionMiddleware))
		return true
	})
	return out
}

// receiverTypeName returns the bare type name of a method receiver, for both
// `(e SysUser)` and `(e *SysUser)`.
func receiverTypeName(expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// callsPackageFunc reports whether body calls name on a package whose import
// path ends in actionsPkgSuffix.
func callsPackageFunc(sf *sourceFile, body ast.Node, name string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != name {
			return true
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		if path, ok := sf.imports[pkg.Name]; ok && strings.HasSuffix(path, actionsPkgSuffix) {
			found = true
			return false
		}
		return true
	})
	return found
}

// apiTypeOf recognises `apis.SysUser{}` and returns the package path and type.
func apiTypeOf(sf *sourceFile, expr ast.Expr) (handlerKey, bool) {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return handlerKey{}, false
	}
	sel, ok := lit.Type.(*ast.SelectorExpr)
	if !ok {
		return handlerKey{}, false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return handlerKey{}, false
	}
	path, ok := sf.imports[pkg.Name]
	if !ok {
		return handlerKey{}, false
	}
	return handlerKey{Pkg: path, Type: sel.Sel.Name}, true
}

// groupSource reports whether expr builds a router group, and names the
// variable it was built from when there is one.
func groupSource(expr ast.Expr) (string, bool) {
	parent := ""
	isGroup := false
	ast.Inspect(expr, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Group" {
			return true
		}
		isGroup = true
		if id, ok := sel.X.(*ast.Ident); ok {
			parent = id.Name
		}
		return true
	})
	return parent, isGroup
}

// groupPath returns the literal path a group was declared with, or "" when it
// is not a literal - a computed prefix is left out of the message rather than
// printed as something it is not.
func groupPath(expr ast.Expr) string {
	out := ""
	ast.Inspect(expr, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Group" || len(call.Args) == 0 {
			return true
		}
		if lit, ok := call.Args[0].(*ast.BasicLit); ok {
			out = strings.Trim(lit.Value, `"`)
		}
		return true
	})
	return out
}

// containsCallNamed reports whether expr contains a call to a function with
// this name, at any depth of a method chain or argument list.
func containsCallNamed(expr ast.Node, name string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fun := call.Fun.(type) {
		case *ast.SelectorExpr:
			if fun.Sel.Name == name {
				found = true
				return false
			}
		case *ast.Ident:
			if fun.Name == name {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// receiverIdentOf returns the variable a `x.method(...)` call was made on.
func receiverIdentOf(call *ast.CallExpr, method string) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != method {
		return "", false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", false
	}
	return id.Name, true
}

// httpMethods are the registration calls this check understands. Any and Match
// are absent on purpose: they take the method as data, and a check that half
// understands a registration is worse than one that says nothing about it.
var httpMethods = map[string]bool{
	"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true, "HEAD": true, "OPTIONS": true,
}

// routeRegistration recognises `g.GET(...)` and names the group and method.
func routeRegistration(call *ast.CallExpr) (string, string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || !httpMethods[sel.Sel.Name] {
		return "", "", false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", "", false
	}
	return id.Name, sel.Sel.Name, true
}

// routeArgs pulls the path and the `api.Handler` argument out of a
// registration, ignoring any middleware written between them.
func routeArgs(call *ast.CallExpr) (route, handlerVar, handlerName string, ok bool) {
	if len(call.Args) < 2 {
		return "", "", "", false
	}
	lit, isLit := call.Args[0].(*ast.BasicLit)
	if !isLit {
		return "", "", "", false
	}
	route = strings.Trim(lit.Value, `"`)
	// The handler is the last argument; anything before it is middleware.
	sel, isSel := call.Args[len(call.Args)-1].(*ast.SelectorExpr)
	if !isSel {
		return "", "", "", false
	}
	id, isIdent := sel.X.(*ast.Ident)
	if !isIdent {
		return "", "", "", false
	}
	return route, id.Name, sel.Sel.Name, true
}
