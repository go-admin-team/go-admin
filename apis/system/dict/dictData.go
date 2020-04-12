package dict

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/app"
	"go-admin/pkg/utils"
	"net/http"
)

// @Summary 字典数据列表
// @Description 获取JSON
// @Tags 字典数据
// @Param status query string false "status"
// @Param dictCode query string false "dictCode"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/data/list [get]
// @Security
func GetDictDataList(c *gin.Context) {
	var data models.DictData
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.DictLabel = c.Request.FormValue("dictLabel")
	data.Status = c.Request.FormValue("status")
	data.DictType = c.Request.FormValue("dictType")
	id := c.Request.FormValue("dictCode")
	data.DictCode, _ = utils.StringToInt(id)
	data.DataScope = utils.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.HasError(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res app.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 通过编码获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param dictCode path int true "字典编码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/data/{dictCode} [get]
// @Security
func GetDictData(c *gin.Context) {
	var DictData models.DictData
	DictData.DictLabel = c.Request.FormValue("dictLabel")
	DictData.DictCode, _ = utils.StringToInt(c.Param("dictCode"))
	result, err := DictData.GetByCode()
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 通过字典类型获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param dictType path int true "dictType"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/databyType/{dictType} [get]
// @Security
func GetDictDataByDictType(c *gin.Context) {
	var DictData models.DictData
	DictData.DictType = c.Param("dictType")
	result, err := DictData.Get()
	pkg.HasError(err, "抱歉未找到相关信息", -1)

	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 添加字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/data [post]
// @Security Bearer
func InsertDictData(c *gin.Context) {
	var data models.DictData
	err := c.BindWith(&data, binding.JSON)
	data.CreateBy = utils.GetUserIdStr(c)
	pkg.HasError(err, "", 500)
	result, err := data.Create()
	pkg.HasError(err, "", -1)
	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 修改字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/data [put]
// @Security Bearer
func UpdateDictData(c *gin.Context) {
	var data models.DictData
	err := c.BindWith(&data, binding.JSON)
	data.UpdateBy = utils.GetUserIdStr(c)
	pkg.HasError(err, "", -1)
	result, err := data.Update(data.DictCode)
	pkg.HasError(err, "", -1)
	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 删除字典数据
// @Description 删除数据
// @Tags 字典数据
// @Param dictCode path int true "dictCode"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dict/data/{dictCode} [delete]
func DeleteDictData(c *gin.Context) {
	var data models.DictData
	id, err := utils.StringToInt(c.Param("dictCode"))
	data.UpdateBy = utils.GetUserIdStr(c)
	_, err = data.Delete(id)
	pkg.HasError(err, "修改失败", 500)

	var res app.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}
