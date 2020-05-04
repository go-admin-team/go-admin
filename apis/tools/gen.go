package tools

import (
	"bytes"
	"github.com/gin-gonic/gin"
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
	tab, _ := table.Get()
	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)

	mp := make(map[string]interface{})
	mp["template/model.go.template"] = b1.String()
	mp["template/api.go.template"] = b2.String()
	mp["template/js.go.template"] = b3.String()
	mp["template/vue.go.template"] = b4.String()
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
	t3, err := template.ParseFiles("template/js.go.template")
	tools2.HasError(err, "", -1)
	t4, err := template.ParseFiles("template/vue.go.template")
	tools2.HasError(err, "", -1)
	tab, _ := table.Get()
	filename := tools2.GetCurrntTimeStr2()
	_ = tools2.PathCreate("./temp/go-admin-" + filename + "/go-admin/apis/")
	_ = tools2.PathCreate("./temp/go-admin-" + filename + "/go-admin/apis/")
	_ = tools2.PathCreate("./temp/go-admin-" + filename + "/go-admin/models/")
	_ = tools2.PathCreate("./temp/go-admin-" + filename + "/go-admin-ui/src/views/" + tab.ModuleName + "/")
	_ = tools2.PathCreate("./temp/go-admin-" + filename + "/go-admin-ui/src/api")

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	tools2.FileCreate(b1, "./temp/go-admin-"+filename+"/go-admin/models/"+tab.PackageName+".go")
	tools2.FileCreate(b2, "./temp/go-admin-"+filename+"/go-admin/apis/"+tab.PackageName+".go")
	tools2.FileCreate(b3, "./temp/go-admin-"+filename+"/go-admin-ui/src/api/"+tab.PackageName+".js")
	tools2.FileCreate(b4, "./temp/go-admin-"+filename+"/go-admin-ui/src/views/"+tab.ModuleName+"/index.vue")

	//tools2.
	tools2.FileZip("./temp/go-admin-"+filename+".zip", "./temp/go-admin-"+filename, "temp/go-admin-"+filename)
	tools2.PathRemove("./temp/go-admin-" + filename)
	c.FileAttachment("./temp/go-admin-"+filename+".zip", "go-admin.zip")
	tools2.PathRemove("./temp/go-admin-" + filename + ".zip")
}
