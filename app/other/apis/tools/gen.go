package tools

import (
	"bytes"
	"fmt"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"go-admin/app/other/models/tools"
)

type Gen struct {
	api.Api
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
	t3, err := template.ParseFiles("template/v4/js.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("js模版读取失败！错误详情：%s", err.Error()))
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

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(500, err, fmt.Sprintf("数据库链接获取失败！错误详情：%s", err.Error()))
		return
	}

	tab, _ := table.Get(db,false)
	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)
	var b7 bytes.Buffer
	err = t7.Execute(&b7, tab)

	mp := make(map[string]interface{})
	mp["template/model.go.template"] = b1.String()
	mp["template/api.go.template"] = b2.String()
	mp["template/js.go.template"] = b3.String()
	mp["template/vue.go.template"] = b4.String()
	mp["template/router.go.template"] = b5.String()
	mp["template/dto.go.template"] = b6.String()
	mp["template/service.go.template"] = b7.String()
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
	tab, _ := table.Get(db,false)

	e.NOActionsGen(c, tab)

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
	tab, _ := table.Get(db,false)
	e.genApiToFile(c, tab)

	e.OK("", "Code generated successfully！")
}

func (e Gen) NOActionsGen(c *gin.Context, tab tools.SysTables) {
	e.Context = c
	log := e.GetLogger()
	tab.MLTBName = strings.Replace(tab.TBName, "_", "-", -1)

	basePath := "template/v4/"
	routerFile := basePath + "no_actions/router_check_role.go.template"

	if tab.IsAuth == 2 {
		routerFile = basePath + "no_actions/router_no_check_role.go.template"
	}

	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("model模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t2, err := template.ParseFiles(basePath + "no_actions/apis.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("api模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("路由模版失败！错误详情：%s", err.Error()))
		return
	}
	t4, err := template.ParseFiles(basePath + "js.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("js模版解析失败！错误详情：%s", err.Error()))
		return
	}
	t5, err := template.ParseFiles(basePath + "vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("vue模版解析失败！错误详情：%s", err.Error()))
		return
	}
	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("dto模版解析失败失败！错误详情：%s", err.Error()))
		return
	}
	t7, err := template.ParseFiles(basePath + "no_actions/service.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("service模版失败！错误详情：%s", err.Error()))
		return
	}

	_ = pkg.PathCreate("./app/" + tab.PackageName + "/apis/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/models/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/router/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/service/dto/")
	_ = pkg.PathCreate(config.GenConfig.FrontPath + "/api/" + tab.PackageName + "/")
	err = pkg.PathCreate(config.GenConfig.FrontPath + "/views/" + tab.PackageName + "/" + tab.MLTBName + "/")
	if err != nil {
		log.Error(err)
		e.Error(500, err, fmt.Sprintf("views目录创建失败！错误详情：%s", err.Error()))
		return
	}

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)
	var b7 bytes.Buffer
	err = t7.Execute(&b7, tab)
	pkg.FileCreate(b1, "./app/"+tab.PackageName+"/models/"+tab.TBName+".go")
	pkg.FileCreate(b2, "./app/"+tab.PackageName+"/apis/"+tab.TBName+".go")
	pkg.FileCreate(b3, "./app/"+tab.PackageName+"/router/"+tab.TBName+".go")
	pkg.FileCreate(b4, config.GenConfig.FrontPath+"/api/"+tab.PackageName+"/"+tab.MLTBName+".js")
	pkg.FileCreate(b5, config.GenConfig.FrontPath+"/views/"+tab.PackageName+"/"+tab.MLTBName+"/index.vue")
	pkg.FileCreate(b6, "./app/"+tab.PackageName+"/service/dto/"+tab.TBName+".go")
	pkg.FileCreate(b7, "./app/"+tab.PackageName+"/service/"+tab.TBName+".go")

}

func (e Gen) genApiToFile(c *gin.Context, tab tools.SysTables) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	basePath := "template/"

	t1, err := template.ParseFiles(basePath + "api_migrate.template")
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("数据迁移模版解析失败！错误详情：%s", err.Error()))
		return
	}
	i := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	var b1 bytes.Buffer
	err = t1.Execute(&b1, struct {
		tools.SysTables
		GenerateTime string
	}{tab, i})

	pkg.FileCreate(b1, "./cmd/migrate/migration/version-local/"+i+"_migrate.go")

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
	tab, _ := table.Get(e.Orm,true)
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
