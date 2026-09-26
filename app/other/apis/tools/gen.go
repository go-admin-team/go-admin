package tools

import (
	"bytes"
	"fmt"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"

	"go-admin/app/other/models/tools"
)

type Gen struct {
	api.Api
}

// genLangFuncs backs the lang-zh/lang-en templates (PRD 010 F3/F9). The
// generated files are TypeScript, and go-admin-ui's eslint config requires
// single-quoted strings with no trailing comma (@stylistic/quotes,
// @stylistic/comma-dangle: never) - text/template's builtin `printf "%q"`
// only produces Go/JSON-style double-quoted output, so this supplies a
// single-quote equivalent instead of leaning on the builtin.
var genLangFuncs = template.FuncMap{
	"singleQuote": func(s string) string {
		r := strings.NewReplacer(`\`, `\\`, `'`, `\'`, "\n", `\n`, "\r", `\r`)
		return "'" + r.Replace(s) + "'"
	},
}

// parseGenTemplate is template.ParseFiles plus genLangFuncs, for the two
// language-pack templates. template.New's name must match the file's base
// name - ParseFiles reuses the template already registered under that name
// instead of creating an unnamed second one, which is what makes Execute
// find the parsed content afterwards.
func parseGenTemplate(path string) (*template.Template, error) {
	base := path[strings.LastIndex(path, "/")+1:]
	return template.New(base).Funcs(genLangFuncs).ParseFiles(path)
}

func (e Gen) Preview(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("tableId接收失败！错误详情：%s", err.Error()))
		return
	}
	table.TableId = id
	t1, err := template.ParseFiles("template/v4/model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("model模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t2, err := template.ParseFiles("template/v4/no_actions/apis.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("api模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t3, err := template.ParseFiles("template/v4/ts.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("ts模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t4, err := template.ParseFiles("template/v4/vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("vue模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t5, err := template.ParseFiles("template/v4/no_actions/router_check_role.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("路由模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t6, err := template.ParseFiles("template/v4/dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("dto模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t7, err := template.ParseFiles("template/v4/no_actions/service.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("service模版读取失败！错误详情：%s", err.Error()))
		return
	}
	// t8/t9 back F3/F9 (PRD 010): one language pack per locale, nested under
	// gen/{PackageName}/{BusinessName}.ts by NOActionsGen below so go-admin-ui's
	// gen-namespace.ts glob (`./*/*.ts` under each locale's gen/) picks them up.
	// See docs-prd/010-代码生成器前端模板迁移Vue3/API契约.md §2.3.
	t8, err := parseGenTemplate("template/v4/lang-zh.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("zh语言包模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t9, err := parseGenTemplate("template/v4/lang-en.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("en语言包模版读取失败！错误详情：%s", err.Error()))
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(500, err, fmt.Sprintf("数据库链接获取失败！错误详情：%s", err.Error()))
		return
	}

	tab, _ := table.Get(db, false)
	// MLTBName (table_name with underscores turned to dashes) is a gorm:"-"
	// field - table.Get never fills it in, so every template that reads it
	// (the .vue/.ts import paths, e.g. "@/api/{PackageName}/{MLTBName}")
	// silently rendered it empty here. NOActionsGen has set this since it
	// existed (see below); Preview never did, which is why the two paths
	// are not interchangeable stand-ins for each other and should not be
	// assumed to be.
	tab.MLTBName = strings.Replace(tab.TBName, "_", "-", -1)
	if err := requireSinglePrimaryKey(tab); err != nil {
		e.Error(500, err, err.Error())
		return
	}
	// R2: infer a width for any column the config page left at colWidth's 0
	// sentinel, before vue.go.template reads .ColWidth - see column_width.go.
	applyInferredColumnWidths(tab.Columns)
	out, err := renderAll(tab, t1, t2, t3, t4, t5, t6, t7, t8, t9)
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("模版渲染失败！错误详情：%s", err.Error()))
		return
	}

	mp := make(map[string]interface{})
	mp["template/model.go.template"] = string(out[0])
	mp["template/api.go.template"] = string(out[1])
	mp["template/api.ts.template"] = string(out[2])
	mp["template/vue.go.template"] = string(out[3])
	mp["template/router.go.template"] = string(out[4])
	mp["template/dto.go.template"] = string(out[5])
	mp["template/service.go.template"] = string(out[6])
	mp["template/lang-zh.go.template"] = string(out[7])
	mp["template/lang-en.go.template"] = string(out[8])
	e.OK(mp, "")
}

func (e Gen) GenCode(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("tableId参数接收失败！错误详情：%s", err.Error()))
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(500, err, fmt.Sprintf("数据库链接获取失败！错误详情：%s", err.Error()))
		return
	}

	table.TableId = id
	tab, _ := table.Get(db, false)

	if !e.NOActionsGen(c, tab) {
		return
	}

	e.OK("", "Code generated successfully！")
}

func (e Gen) GenApiToFile(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("tableId参数获取失败！错误详情：%s", err.Error()))
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(500, err, fmt.Sprintf("数据库链接获取失败！错误详情：%s", err.Error()))
		return
	}

	table.TableId = id
	tab, _ := table.Get(db, false)
	if !e.genApiToFile(c, tab) {
		return
	}

	e.OK("", "Code generated successfully！")
}

// NOActionsGen renders every template for tab and writes the files, and
// reports whether it did. On false it has already written the error response,
// and nothing has been written to disk if a template failed to render.
func (e Gen) NOActionsGen(c *gin.Context, tab tools.SysTables) bool {
	e.Context = c
	log := e.GetLogger()
	tab.MLTBName = strings.Replace(tab.TBName, "_", "-", -1)
	if err := requireSinglePrimaryKey(tab); err != nil {
		e.Error(500, err, err.Error())
		return false
	}
	// R2: see the matching call and comment in Preview above.
	applyInferredColumnWidths(tab.Columns)

	basePath := "template/v4/"
	routerFile := basePath + "no_actions/router_check_role.go.template"

	if tab.IsAuth == 2 {
		routerFile = basePath + "no_actions/router_no_check_role.go.template"
	}

	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("model模版读取失败！错误详情：%s", err.Error()))
		return false
	}
	t2, err := template.ParseFiles(basePath + "no_actions/apis.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("api模版读取失败！错误详情：%s", err.Error()))
		return false
	}
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("路由模版失败！错误详情：%s", err.Error()))
		return false
	}
	t4, err := template.ParseFiles(basePath + "ts.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("ts模版解析失败！错误详情：%s", err.Error()))
		return false
	}
	t5, err := template.ParseFiles(basePath + "vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("vue模版解析失败！错误详情：%s", err.Error()))
		return false
	}
	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("dto模版解析失败失败！错误详情：%s", err.Error()))
		return false
	}
	t7, err := template.ParseFiles(basePath + "no_actions/service.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("service模版失败！错误详情：%s", err.Error()))
		return false
	}
	// t8/t9 back F3/F9 (PRD 010): see the matching comment in Preview above.
	t8, err := parseGenTemplate(basePath + "lang-zh.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("zh语言包模版解析失败！错误详情：%s", err.Error()))
		return false
	}
	t9, err := parseGenTemplate(basePath + "lang-en.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("en语言包模版解析失败！错误详情：%s", err.Error()))
		return false
	}

	out, err := renderAll(tab, t1, t2, t3, t4, t5, t6, t7, t8, t9)
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("模版渲染失败！错误详情：%s", err.Error()))
		return false
	}

	// gen/{PackageName}/ nests under each locale so go-admin-ui's
	// gen-namespace.ts (`./*/*.ts` glob, one level under gen/) picks the file
	// up - a flat gen/{BusinessName}.ts would let two tables in different
	// packages silently overwrite each other's translations, since
	// BusinessName only has a pattern check, no uniqueness check.
	files := []struct {
		path    string
		content []byte
	}{
		{"./app/" + tab.PackageName + "/models/" + tab.TBName + ".go", out[0]},
		{"./app/" + tab.PackageName + "/apis/" + tab.TBName + ".go", out[1]},
		{"./app/" + tab.PackageName + "/router/" + tab.TBName + ".go", out[2]},
		{config.GenConfig.FrontPath + "/api/" + tab.PackageName + "/" + tab.MLTBName + ".ts", out[3]},
		{config.GenConfig.FrontPath + "/views/" + tab.PackageName + "/" + tab.MLTBName + "/index.vue", out[4]},
		{"./app/" + tab.PackageName + "/service/dto/" + tab.TBName + ".go", out[5]},
		{"./app/" + tab.PackageName + "/service/" + tab.TBName + ".go", out[6]},
		{config.GenConfig.FrontPath + "/lang/zh-CN/gen/" + tab.PackageName + "/" + tab.BusinessName + ".ts", out[7]},
		{config.GenConfig.FrontPath + "/lang/en-US/gen/" + tab.PackageName + "/" + tab.BusinessName + ".ts", out[8]},
	}
	for _, f := range files {
		if err := writeGenerated(f.path, f.content); err != nil {
			log.Error(err)
			e.Error(500, err, fmt.Sprintf("生成文件写入失败！错误详情：%s", err.Error()))
			return false
		}
	}
	return true
}

// genApiToFile writes the migration that seeds tab's menu and APIs, and
// reports whether it did; on false the error response is already written.
func (e Gen) genApiToFile(c *gin.Context, tab tools.SysTables) bool {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return false
	}

	basePath := "template/"

	t1, err := template.ParseFiles(basePath + "api_migrate.template")
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("数据迁移模版解析失败！错误详情：%s", err.Error()))
		return false
	}
	i := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	out, err := renderAll(struct {
		tools.SysTables
		GenerateTime string
	}{tab, i}, t1)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("数据迁移模版渲染失败！错误详情：%s", err.Error()))
		return false
	}
	if err := writeGenerated("./cmd/migrate/migration/version-local/"+i+"_migrate.go", out[0]); err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("数据迁移文件写入失败！错误详情：%s", err.Error()))
		return false
	}
	return true
}

func (e Gen) GenMenuAndApi(c *gin.Context) {
	s := service.SysMenu{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("tableId参数解析失败！错误详情：%s", err.Error()))
		return
	}

	table.TableId = id
	tab, _ := table.Get(e.Orm, true)
	tab.MLTBName = strings.Replace(tab.TBName, "_", "-", -1)

	Mmenu := dto.SysMenuInsertReq{}
	Mmenu.Title = tab.TableComment
	Mmenu.Icon = "pass"
	Mmenu.Path = "/" + tab.MLTBName
	Mmenu.MenuType = "M"
	Mmenu.Action = "无"
	Mmenu.ParentId = 0
	Mmenu.NoCache = false
	Mmenu.Component = "Layout"
	Mmenu.Sort = 0
	Mmenu.Visible = "0"
	Mmenu.IsFrame = "0"
	Mmenu.CreateBy = 1
	s.Insert(&Mmenu)

	Cmenu := dto.SysMenuInsertReq{}
	Cmenu.MenuName = tab.ClassName + "Manage"
	Cmenu.Title = tab.TableComment
	Cmenu.Icon = "pass"
	Cmenu.Path = "/" + tab.PackageName + "/" + tab.MLTBName
	Cmenu.MenuType = "C"
	Cmenu.Action = "无"
	Cmenu.Permission = tab.PackageName + ":" + tab.BusinessName + ":list"
	Cmenu.ParentId = Mmenu.MenuId
	Cmenu.NoCache = false
	Cmenu.Component = "/" + tab.PackageName + "/" + tab.MLTBName + "/index"
	Cmenu.Sort = 0
	Cmenu.Visible = "0"
	Cmenu.IsFrame = "0"
	Cmenu.CreateBy = 1
	Cmenu.UpdateBy = 1
	s.Insert(&Cmenu)

	MList := dto.SysMenuInsertReq{}
	MList.MenuName = ""
	MList.Title = "分页获取" + tab.TableComment
	MList.Icon = ""
	MList.Path = tab.TBName
	MList.MenuType = "F"
	MList.Action = "无"
	MList.Permission = tab.PackageName + ":" + tab.BusinessName + ":query"
	MList.ParentId = Cmenu.MenuId
	MList.NoCache = false
	MList.Sort = 0
	MList.Visible = "0"
	MList.IsFrame = "0"
	MList.CreateBy = 1
	MList.UpdateBy = 1
	s.Insert(&MList)

	MCreate := dto.SysMenuInsertReq{}
	MCreate.MenuName = ""
	MCreate.Title = "创建" + tab.TableComment
	MCreate.Icon = ""
	MCreate.Path = tab.TBName
	MCreate.MenuType = "F"
	MCreate.Action = "无"
	MCreate.Permission = tab.PackageName + ":" + tab.BusinessName + ":add"
	MCreate.ParentId = Cmenu.MenuId
	MCreate.NoCache = false
	MCreate.Sort = 0
	MCreate.Visible = "0"
	MCreate.IsFrame = "0"
	MCreate.CreateBy = 1
	MCreate.UpdateBy = 1
	s.Insert(&MCreate)

	MUpdate := dto.SysMenuInsertReq{}
	MUpdate.MenuName = ""
	MUpdate.Title = "修改" + tab.TableComment
	MUpdate.Icon = ""
	MUpdate.Path = tab.TBName
	MUpdate.MenuType = "F"
	MUpdate.Action = "无"
	MUpdate.Permission = tab.PackageName + ":" + tab.BusinessName + ":edit"
	MUpdate.ParentId = Cmenu.MenuId
	MUpdate.NoCache = false
	MUpdate.Sort = 0
	MUpdate.Visible = "0"
	MUpdate.IsFrame = "0"
	MUpdate.CreateBy = 1
	MUpdate.UpdateBy = 1
	s.Insert(&MUpdate)

	MDelete := dto.SysMenuInsertReq{}
	MDelete.MenuName = ""
	MDelete.Title = "删除" + tab.TableComment
	MDelete.Icon = ""
	MDelete.Path = tab.TBName
	MDelete.MenuType = "F"
	MDelete.Action = "无"
	MDelete.Permission = tab.PackageName + ":" + tab.BusinessName + ":remove"
	MDelete.ParentId = Cmenu.MenuId
	MDelete.NoCache = false
	MDelete.Sort = 0
	MDelete.Visible = "0"
	MDelete.IsFrame = "0"
	MDelete.CreateBy = 1
	MDelete.UpdateBy = 1
	s.Insert(&MDelete)

	e.OK("", "数据生成成功！")
}

// requireSinglePrimaryKey refuses a table whose primary key is not exactly one
// column. The templates address a row by a single key - one field on the
// model, one path parameter, one value per row in a delete - so a table with
// no primary key, or a composite one, has no shape they can render: generating
// it anyway writes a model whose GetId names nothing, or silently keys the
// whole table on whichever key column happened to be imported last.
func requireSinglePrimaryKey(tab tools.SysTables) error {
	var keys []string
	for _, c := range tab.Columns {
		if c.Pk {
			keys = append(keys, c.ColumnName)
		}
	}
	switch len(keys) {
	case 1:
		return nil
	case 0:
		return fmt.Errorf("表 %s 没有主键，代码生成需要恰好一个主键列", tab.TBName)
	default:
		return fmt.Errorf("表 %s 是联合主键（%s），代码生成需要恰好一个主键列", tab.TBName, strings.Join(keys, ", "))
	}
}

// renderAll executes each template against data and returns the outputs in
// the same order, stopping at the first template that fails. The caller
// writes nothing until every template has rendered: a template that fails
// part-way leaves a truncated file behind it, which compiles as nothing and
// at a glance looks like the real thing.
func renderAll(data any, tpls ...*template.Template) ([][]byte, error) {
	out := make([][]byte, len(tpls))
	for i, t := range tpls {
		var b bytes.Buffer
		if err := t.Execute(&b, data); err != nil {
			return nil, err
		}
		out[i] = b.Bytes()
	}
	return out, nil
}

// writeGenerated writes one generated file, creating its directory first.
// It stands in for pkg.FileCreate, which returns no error at all and, when
// the file cannot be created, closes a nil file and ends the process with
// log.Fatalln - one unwritable path took the whole server down with it.
func writeGenerated(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}
