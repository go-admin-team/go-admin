package sysfile

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"go-admin/app/admin/models"
	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/app/msg"
)

func GetSysFileInfoOldList(c *gin.Context) {
	var data models.SysFileInfoOld
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = tools.StringToInt(size)
	}
	tools.HasError(err, "", -1)
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = tools.StringToInt(index)
	}
	tools.HasError(err, "", -1)
	if pid := c.Request.FormValue("pId"); pid != "" {
		data.PId, err = tools.StringToInt(pid)
	}
	tools.HasError(err, "", -1)

	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

func GetSysFileInfoOld(c *gin.Context) {
	var data models.SysFileInfoOld
	data.Id, _ = tools.StringToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

func InsertSysFileInfoOld(c *gin.Context) {
	var data models.SysFileInfoOld
	err := c.ShouldBindJSON(&data)
	data.CreateBy = tools.GetUserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

func UpdateSysFileInfoOld(c *gin.Context) {
	var data models.SysFileInfoOld
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

func DeleteSysFileInfoOld(c *gin.Context) {
	var data models.SysFileInfoOld
	data.UpdateBy = tools.GetUserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup("id", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
