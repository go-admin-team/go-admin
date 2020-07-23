package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/app/msg"
	"net/http"
)

// @Summary 配置列表数据
// @Description 获取JSON
// @Tags 配置
// @Param configKey query string false "configKey"
// @Param configName query string false "configName"
// @Param configType query string false "configType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/configList [get]
// @Security Bearer
func GetConfigList(c *gin.Context) {
	var data models.SysConfig
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	data.ConfigKey = c.Request.FormValue("configKey")
	data.ConfigName = c.Request.FormValue("configName")
	data.ConfigType = c.Request.FormValue("configType")
	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res app.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 配置
// @Param configId path int true "配置编码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/config/{configId} [get]
// @Security Bearer
func GetConfig(c *gin.Context) {
	var Config models.SysConfig
	Config.ConfigId, _ = tools.StringToInt(c.Param("configId"))
	result, err := Config.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	var res app.Response
	res.Data = result

	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 配置
// @Param configKey path int true "configKey"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/configKey/{configKey} [get]
// @Security Bearer
func GetConfigByConfigKey(c *gin.Context) {
	var Config models.SysConfig
	Config.ConfigKey = c.Param("configKey")
	result, err := Config.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, result.ConfigValue)
}

// @Summary 添加配置
// @Description 获取JSON
// @Tags 配置
// @Accept  application/json
// @Product application/json
// @Param data body models.SysConfig true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/data [post]
// @Security Bearer
func InsertConfig(c *gin.Context) {
	var data models.SysConfig
	err := c.BindWith(&data, binding.JSON)
	data.CreateBy = tools.GetUserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

// @Summary 修改配置
// @Description 获取JSON
// @Tags 配置
// @Accept  application/json
// @Product application/json
// @Param data body models.SysConfig true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/config [put]
// @Security Bearer
func UpdateConfig(c *gin.Context) {
	var data models.SysConfig
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.ConfigId)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

// @Summary 删除配置
// @Description 删除数据
// @Tags 配置
// @Param configId path int true "configId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/config/{configId} [delete]
func DeleteConfig(c *gin.Context) {
	var data models.SysConfig
	data.UpdateBy = tools.GetUserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup("configId", c)
	result, err := data.BatchDelete(IDS)
	tools.HasError(err, "修改失败", 500)
	app.OK(c, result, msg.DeletedSuccess)
}
