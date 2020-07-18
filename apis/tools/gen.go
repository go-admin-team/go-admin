package tools

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/tools"
	tools2 "go-admin/tools"
	"go-admin/tools/app"
	"net/http"
	"text/template"
)

func Preview(c *gin.Context) {
	table := tools.SysTables{}
	id, err := tools2.StringToInt(c.Param("tableId"))
	tools2.HasError(err, "", -1)
	table.TableId = id
	t1, err := template.ParseFiles("template/model.go.template")
	tools2.HasError(err, "", -1)
	t2, err := template.ParseFiles("template/api.go.template")
	tools2.HasError(err, "", -1)
	t3, err := template.ParseFiles("template/js.go.template")
	tools2.HasError(err, "", -1)
	t4, err := template.ParseFiles("template/vue.go.template")
	tools2.HasError(err, "", -1)
	t5, err := template.ParseFiles("template/router.go.template")
	tools2.HasError(err, "", -1)
	tab, _ := table.Get()
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

	mp := make(map[string]interface{})
	mp["template/model.go.template"] = b1.String()
	mp["template/api.go.template"] = b2.String()
	mp["template/js.go.template"] = b3.String()
	mp["template/vue.go.template"] = b4.String()
	mp["template/router.go.template"] = b5.String()
	var res app.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

func GenCode(c *gin.Context) {
	table := tools.SysTables{}
	id, err := tools2.StringToInt(c.Param("tableId"))
	tools2.HasError(err, "", -1)
	table.TableId = id

	t1, err := template.ParseFiles("template/model.go.template")
	tools2.HasError(err, "", -1)
	t2, err := template.ParseFiles("template/api.go.template")
	tools2.HasError(err, "", -1)
	t3, err := template.ParseFiles("template/router.go.template")
	tools2.HasError(err, "", -1)
	tab, _ := table.Get()
	_ = tools2.PathCreate("./apis/" + tab.ModuleName + "/")
	_ = tools2.PathCreate("./models/")
	_ = tools2.PathCreate("./router/")

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	tools2.FileCreate(b1, "./models/"+tab.PackageName+".go")
	tools2.FileCreate(b2, "./apis/"+tab.ModuleName+"/"+tab.PackageName+".go")
	tools2.FileCreate(b3, "./router/"+tab.PackageName+".go")
	app.OK(c, "", "代码生成成功！")
}

func GenMenuAndApi(c *gin.Context) {

	table := tools.SysTables{}
	timeNow := tools2.GetCurrentTime()
	id, err := tools2.StringToInt(c.Param("tableId"))
	tools2.HasError(err, "", -1)
	table.TableId = id
	tab, _ := table.Get()
	Mmenu := models.Menu{}
	Mmenu.MenuName = tab.TBName + "管理"
	Mmenu.Title = tab.TableComment
	Mmenu.Icon = "pass"
	Mmenu.Path = tab.TBName
	Mmenu.MenuType = "M"
	Mmenu.Action = "无"
	Mmenu.Permission = tab.PackageName + ":" + tab.ModuleName + ":list"
	Mmenu.ParentId = 0
	Mmenu.NoCache = false
	Mmenu.Component = "/" + tab.ModuleName + "/index"
	Mmenu.Sort = 0
	Mmenu.Visible = "0"
	Mmenu.IsFrame = "0"
	Mmenu.CreateBy = "1"
	Mmenu.UpdateBy = "1"
	Mmenu.CreatedAt = timeNow
	Mmenu.UpdatedAt = timeNow
	Mmenu.MenuId, err = Mmenu.Create()

	MList := models.Menu{}
	MList.MenuName = "分页获取" + tab.TBName
	MList.Title = tab.TableComment
	MList.Icon = "pass"
	MList.Path = tab.TBName
	MList.MenuType = "F"
	MList.Action = "无"
	MList.Permission = tab.PackageName + ":" + tab.ModuleName + ":query"
	MList.ParentId = Mmenu.MenuId
	MList.NoCache = false
	MList.Sort = 0
	MList.Visible = "0"
	MList.IsFrame = "0"
	MList.CreateBy = "1"
	MList.UpdateBy = "1"
	MList.CreatedAt = timeNow
	MList.UpdatedAt = timeNow
	MList.MenuId, err = MList.Create()

	MCreate := models.Menu{}
	MCreate.MenuName = tab.TBName
	MCreate.Title = "创建" + tab.TableComment
	MCreate.Icon = "pass"
	MCreate.Path = tab.TBName
	MCreate.MenuType = "F"
	MCreate.Action = "无"
	MCreate.Permission = tab.PackageName + ":" + tab.ModuleName + ":add"
	MCreate.ParentId = Mmenu.MenuId
	MCreate.NoCache = false
	MCreate.Sort = 0
	MCreate.Visible = "0"
	MCreate.IsFrame = "0"
	MCreate.CreateBy = "1"
	MCreate.UpdateBy = "1"
	MCreate.CreatedAt = timeNow
	MCreate.UpdatedAt = timeNow
	MCreate.MenuId, err = MCreate.Create()

	MUpdate := models.Menu{}
	MUpdate.MenuName = tab.TBName
	MUpdate.Title = tab.TableComment
	MUpdate.Icon = "pass"
	MUpdate.Path = "修改" + tab.TBName
	MUpdate.MenuType = "F"
	MUpdate.Action = "无"
	MUpdate.Permission = tab.PackageName + ":" + tab.ModuleName + ":edit"
	MUpdate.ParentId = Mmenu.MenuId
	MUpdate.NoCache = false
	MUpdate.Sort = 0
	MUpdate.Visible = "0"
	MUpdate.IsFrame = "0"
	MUpdate.CreateBy = "1"
	MUpdate.UpdateBy = "1"
	MUpdate.CreatedAt = timeNow
	MUpdate.UpdatedAt = timeNow
	MUpdate.MenuId, err = MUpdate.Create()

	MDelete := models.Menu{}
	MDelete.MenuName = tab.TBName
	MDelete.Title = "删除" + tab.TableComment
	MDelete.Icon = "pass"
	MDelete.Path = tab.TBName
	MDelete.MenuType = "F"
	MDelete.Action = "无"
	MDelete.Permission = tab.PackageName + ":" + tab.ModuleName + ":remove"
	MDelete.ParentId = Mmenu.MenuId
	MDelete.NoCache = false
	MDelete.Sort = 0
	MDelete.Visible = "0"
	MDelete.IsFrame = "0"
	MDelete.CreateBy = "1"
	MDelete.UpdateBy = "1"
	MDelete.CreatedAt = timeNow
	MDelete.UpdatedAt = timeNow
	MDelete.MenuId, err = MDelete.Create()

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
	Amenu.MenuId, err = Amenu.Create()

	AList := models.Menu{}
	AList.MenuName = tab.TBName
	AList.Title = "分页获取" + tab.TableComment
	AList.Icon = "bug"
	AList.Path = "/api/v1/" + tab.ModuleName + "List"
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
	AList.MenuId, err = AList.Create()

	AGet := models.Menu{}
	AGet.MenuName = tab.TBName
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
	AGet.MenuId, err = AGet.Create()

	ACreate := models.Menu{}
	ACreate.MenuName = tab.TBName
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
	ACreate.MenuId, err = ACreate.Create()

	AUpdate := models.Menu{}
	AUpdate.MenuName = tab.TBName
	AUpdate.Title = "修改" + tab.TableComment
	AUpdate.Icon = "bug"
	AUpdate.Path = "/api/v1/" + tab.ModuleName
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
	AUpdate.MenuId, err = AUpdate.Create()

	ADelete := models.Menu{}
	ADelete.MenuName = tab.TBName
	ADelete.Title = "删除" + tab.TableComment
	ADelete.Icon = "bug"
	ADelete.Path = "/api/v1/" + tab.ModuleName + "/:id"
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
	ADelete.MenuId, err = ADelete.Create()

	app.OK(c, "", "数据生成成功！")
}
