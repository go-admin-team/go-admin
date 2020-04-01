package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/utils"
	"net/http"
	"strconv"
)

// @Summary 登录日志列表
// @Description 获取JSON
// @Tags 登录日志
// @Param status query string false "status"
// @Param dictCode query string false "dictCode"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/loginloglist [get]
// @Security
func GetLoginLogList(c *gin.Context) {
	var data models.LoginLog
	var err error
	var pageSize = 10
	var pageIndex = 1

	size := c.Request.FormValue("pageSize")
	if size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	index := c.Request.FormValue("pageIndex")
	if index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.UserName = c.Request.FormValue("userName")
	data.Status = c.Request.FormValue("status")
	data.Ipaddr = c.Request.FormValue("ipaddr")
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.AssertErr(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res models.Response
	res.Data = mp
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 通过编码获取登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Param infoId path int true "infoId"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/loginlog/{infoId} [get]
// @Security
func GetLoginLog(c *gin.Context) {
	var LoginLog models.LoginLog
	LoginLog.InfoId, _ = strconv.ParseInt(c.Param("infoId"), 10, 64)
	result, err := LoginLog.Get()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)

	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 添加登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Accept  application/json
// @Product application/json
// @Param data body models.LoginLog true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/loginlog [post]
// @Security Bearer
func InsertLoginLog(c *gin.Context) {
	var data models.LoginLog
	err := c.BindWith(&data, binding.JSON)
	pkg.AssertErr(err, "", 500)
	result, err := data.Create()
	pkg.AssertErr(err, "", -1)
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 修改登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Accept  application/json
// @Product application/json
// @Param data body models.LoginLog true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/loginlog [put]
// @Security Bearer
func UpdateLoginLog(c *gin.Context) {
	var data models.LoginLog
	err := c.BindWith(&data, binding.JSON)
	pkg.AssertErr(err, "", -1)
	result, err := data.Update(data.InfoId)
	pkg.AssertErr(err, "", -1)
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 批量删除登录日志
// @Description 删除数据
// @Tags 登录日志
// @Param infoId path string true "以逗号（,）分割的infoId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/loginlog/{infoId} [delete]
func DeleteLoginLog(c *gin.Context) {
	var data models.LoginLog
	data.UpdateBy = utils.GetUserIdStr(c)
	IDS := utils.IdsStrToIdsInt64Group("infoId", c)
	_, err := data.BatchDelete(IDS)
	pkg.AssertErr(err, "修改失败", 500)
	var res models.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}
