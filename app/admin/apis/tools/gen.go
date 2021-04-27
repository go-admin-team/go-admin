package tools

import (
	"bytes"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/tools"
	"go-admin/common/apis"
)

type Gen struct {
	apis.Api
}

func (e *Gen) Preview(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	table.TableId = id
	t1, err := template.ParseFiles("template/v4/model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t2, err := template.ParseFiles("template/v4/no_actions/apis.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t3, err := template.ParseFiles("template/v4/js.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t4, err := template.ParseFiles("template/v4/vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t5, err := template.ParseFiles("template/v4/no_actions/router_check_role.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t6, err := template.ParseFiles("template/v4/dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t7, err := template.ParseFiles("template/v4/no_actions/service.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		return
	}

	tab, _ := table.Get(db)
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

func (e *Gen) GenCodeV3(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		return
	}

	table.TableId = id
	tab, _ := table.Get(db)

	if tab.IsActions == 1 {
		e.ActionsGenV3(c, tab)
	} else {
		e.NOActionsGenV3(c, tab)
	}

	e.OK("", "Code generated successfully！")
}

func (e *Gen) GenApiToFile(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		return
	}

	table.TableId = id
	tab, _ := table.Get(db)
	e.genApiToFile(c, tab)

	e.OK("", "Code generated successfully！")
}

func (e *Gen) NOActionsGenV3(c *gin.Context, tab tools.SysTables) {
	e.Context = c
	log := e.GetLogger()

	basePath := "template/v4/"
	routerFile := basePath + "no_actions/router_check_role.go.template"

	if tab.IsAuth == 2 {
		routerFile = basePath + "no_actions/router_no_check_role.go.template"
	}

	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t2, err := template.ParseFiles(basePath + "no_actions/apis.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t4, err := template.ParseFiles(basePath + "js.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t5, err := template.ParseFiles(basePath + "vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t7, err := template.ParseFiles(basePath + "no_actions/service.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}

	_ = pkg.PathCreate("./app/" + tab.PackageName + "/apis/" + tab.ModuleName)
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/models/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/router/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/service/dto/")
	_ = pkg.PathCreate(config.GenConfig.FrontPath + "/api/")
	_ = pkg.PathCreate(config.GenConfig.FrontPath + "/views/" + tab.BusinessName)

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
	pkg.FileCreate(b1, "./app/"+tab.PackageName+"/models/"+tab.BusinessName+".go")
	pkg.FileCreate(b2, "./app/"+tab.PackageName+"/apis/"+tab.ModuleName+"/"+tab.BusinessName+".go")
	pkg.FileCreate(b3, "./app/"+tab.PackageName+"/router/"+tab.BusinessName+".go")
	pkg.FileCreate(b4, config.GenConfig.FrontPath+"/api/"+tab.BusinessName+".js")
	pkg.FileCreate(b5, config.GenConfig.FrontPath+"/views/"+tab.BusinessName+"/index.vue")
	pkg.FileCreate(b6, "./app/"+tab.PackageName+"/service/dto/"+tab.BusinessName+".go")
	pkg.FileCreate(b7, "./app/"+tab.PackageName+"/service/"+tab.BusinessName+".go")

}

func (e *Gen) genApiToFile(c *gin.Context, tab tools.SysTables) {
	e.Context = c
	log := e.GetLogger()

	basePath := "template/"

	t1, err := template.ParseFiles(basePath + "api_migrate.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
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

func (e *Gen) ActionsGenV3(c *gin.Context, tab tools.SysTables) {
	log := api.GetRequestLogger(c)
	basePath := "template/v4/"
	routerFile := basePath + "actions/router_check_role.go.template"

	if tab.IsAuth == 2 {
		routerFile = basePath + "actions/router_no_check_role.go.template"
	}

	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t4, err := template.ParseFiles(basePath + "js.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t5, err := template.ParseFiles(basePath + "vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}
	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, err, "")
		return
	}

	_ = pkg.PathCreate("./app/" + tab.PackageName + "/models/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/router/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/service/dto/")
	_ = pkg.PathCreate(config.GenConfig.FrontPath + "/api/")
	_ = pkg.PathCreate(config.GenConfig.FrontPath + "/views/" + tab.BusinessName)

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)

	pkg.FileCreate(b1, "./app/"+tab.PackageName+"/models/"+tab.BusinessName+".go")
	pkg.FileCreate(b3, "./app/"+tab.PackageName+"/router/"+tab.BusinessName+".go")
	pkg.FileCreate(b4, config.GenConfig.FrontPath+"/api/"+tab.BusinessName+".js")
	pkg.FileCreate(b5, config.GenConfig.FrontPath+"/views/"+tab.BusinessName+"/index.vue")
	pkg.FileCreate(b6, "./app/"+tab.PackageName+"/service/dto/"+tab.BusinessName+".go")
}

func (e *Gen) GenMenuAndApi(c *gin.Context) {
	log := api.GetRequestLogger(c)

	table := tools.SysTables{}
	timeNow := pkg.GetCurrentTime()
	id, err := pkg.StringToInt(c.Param("tableId"))
	pkg.HasError(err, "", -1)

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		return
	}

	table.TableId = id
	tab, _ := table.Get(db)
	Mmenu := models.Menu{}
	Mmenu.MenuName = tab.TBName + "Manage"
	Mmenu.Title = tab.TableComment
	Mmenu.Icon = "pass"
	Mmenu.Path = "/" + tab.TBName
	Mmenu.MenuType = "M"
	Mmenu.Action = "无"
	Mmenu.ParentId = 0
	Mmenu.NoCache = false
	Mmenu.Component = "Layout"
	Mmenu.Sort = 0
	Mmenu.Visible = "0"
	Mmenu.IsFrame = "0"
	Mmenu.CreateBy = "1"
	Mmenu.UpdateBy = "1"
	Mmenu.CreatedAt = timeNow
	Mmenu.UpdatedAt = timeNow
	Mmenu.MenuId, err = Mmenu.Create(db)

	Cmenu := models.Menu{}
	Cmenu.MenuName = tab.TBName
	Cmenu.Title = tab.TableComment
	Cmenu.Icon = "pass"
	Cmenu.Path = tab.TBName
	Cmenu.MenuType = "C"
	Cmenu.Action = "无"
	Cmenu.Permission = tab.PackageName + ":" + tab.BusinessName + ":list"
	Cmenu.ParentId = Mmenu.MenuId
	Cmenu.NoCache = false
	Cmenu.Component = "/" + tab.BusinessName + "/index"
	Cmenu.Sort = 0
	Cmenu.Visible = "0"
	Cmenu.IsFrame = "0"
	Cmenu.CreateBy = "1"
	Cmenu.UpdateBy = "1"
	Cmenu.CreatedAt = timeNow
	Cmenu.UpdatedAt = timeNow
	Cmenu.MenuId, err = Cmenu.Create(db)

	MList := models.Menu{}
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
	MList.CreateBy = "1"
	MList.UpdateBy = "1"
	MList.CreatedAt = timeNow
	MList.UpdatedAt = timeNow
	MList.MenuId, err = MList.Create(db)

	MCreate := models.Menu{}
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
	MCreate.CreateBy = "1"
	MCreate.UpdateBy = "1"
	MCreate.CreatedAt = timeNow
	MCreate.UpdatedAt = timeNow
	MCreate.MenuId, err = MCreate.Create(db)

	MUpdate := models.Menu{}
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
	MUpdate.CreateBy = "1"
	MUpdate.UpdateBy = "1"
	MUpdate.CreatedAt = timeNow
	MUpdate.UpdatedAt = timeNow
	MUpdate.MenuId, err = MUpdate.Create(db)

	MDelete := models.Menu{}
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
	MDelete.CreateBy = "1"
	MDelete.UpdateBy = "1"
	MDelete.CreatedAt = timeNow
	MDelete.UpdatedAt = timeNow
	MDelete.MenuId, err = MDelete.Create(db)

	var InterfaceId = 63
	Amenu := models.Menu{}
	Amenu.MenuName = tab.TBName
	Amenu.Title = tab.TableComment
	Amenu.Icon = "bug"
	Amenu.Path = tab.TBName
	Amenu.MenuType = "M"
	Amenu.Action = "无"
	Amenu.ParentId = InterfaceId
	Amenu.NoCache = false
	Amenu.Sort = 0
	Amenu.Visible = "1"
	Amenu.IsFrame = "0"
	Amenu.CreateBy = "1"
	Amenu.UpdateBy = "1"
	Amenu.CreatedAt = timeNow
	Amenu.UpdatedAt = timeNow
	Amenu.MenuId, err = Amenu.Create(db)

	AList := models.Menu{}
	AList.MenuName = ""
	AList.Title = "分页获取" + tab.TableComment
	AList.Icon = "bug"
	AList.Path = "/api/v1/" + tab.ModuleName
	AList.MenuType = "A"
	AList.Action = "GET"
	AList.ParentId = Amenu.MenuId
	AList.NoCache = false
	AList.Sort = 0
	AList.Visible = "1"
	AList.IsFrame = "0"
	AList.CreateBy = "1"
	AList.UpdateBy = "1"
	AList.CreatedAt = timeNow
	AList.UpdatedAt = timeNow
	AList.MenuId, err = AList.Create(db)

	AGet := models.Menu{}
	AGet.MenuName = ""
	AGet.Title = "根据id获取" + tab.TableComment
	AGet.Icon = "bug"
	AGet.Path = "/api/v1/" + tab.ModuleName + "/:id"
	AGet.MenuType = "A"
	AGet.Action = "GET"
	AGet.ParentId = Amenu.MenuId
	AGet.NoCache = false
	AGet.Sort = 0
	AGet.Visible = "1"
	AGet.IsFrame = "0"
	AGet.CreateBy = "1"
	AGet.UpdateBy = "1"
	AGet.CreatedAt = timeNow
	AGet.UpdatedAt = timeNow
	AGet.MenuId, err = AGet.Create(db)

	ACreate := models.Menu{}
	ACreate.MenuName = ""
	ACreate.Title = "创建" + tab.TableComment
	ACreate.Icon = "bug"
	ACreate.Path = "/api/v1/" + tab.ModuleName
	ACreate.MenuType = "A"
	ACreate.Action = "POST"
	ACreate.ParentId = Amenu.MenuId
	ACreate.NoCache = false
	ACreate.Sort = 0
	ACreate.Visible = "1"
	ACreate.IsFrame = "0"
	ACreate.CreateBy = "1"
	ACreate.UpdateBy = "1"
	ACreate.CreatedAt = timeNow
	ACreate.UpdatedAt = timeNow
	ACreate.MenuId, err = ACreate.Create(db)

	AUpdate := models.Menu{}
	AUpdate.MenuName = ""
	AUpdate.Title = "修改" + tab.TableComment
	AUpdate.Icon = "bug"
	AUpdate.Path = "/api/v1/" + tab.ModuleName + "/:id"
	AUpdate.MenuType = "A"
	AUpdate.Action = "PUT"
	AUpdate.ParentId = Amenu.MenuId
	AUpdate.NoCache = false
	AUpdate.Sort = 0
	AUpdate.Visible = "1"
	AUpdate.IsFrame = "0"
	AUpdate.CreateBy = "1"
	AUpdate.UpdateBy = "1"
	AUpdate.CreatedAt = timeNow
	AUpdate.UpdatedAt = timeNow
	AUpdate.MenuId, err = AUpdate.Create(db)

	ADelete := models.Menu{}
	ADelete.MenuName = ""
	ADelete.Title = "删除" + tab.TableComment
	ADelete.Icon = "bug"
	ADelete.Path = "/api/v1/" + tab.ModuleName
	ADelete.MenuType = "A"
	ADelete.Action = "DELETE"
	ADelete.ParentId = Amenu.MenuId
	ADelete.NoCache = false
	ADelete.Sort = 0
	ADelete.Visible = "1"
	ADelete.IsFrame = "0"
	ADelete.CreateBy = "1"
	ADelete.UpdateBy = "1"
	ADelete.CreatedAt = timeNow
	ADelete.UpdatedAt = timeNow
	ADelete.MenuId, err = ADelete.Create(db)

	e.OK("", "数据生成成功！")
}
