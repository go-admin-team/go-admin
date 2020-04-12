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

// @Summary 字典类型列表数据
// @Description 获取JSON
// @Tags 字典类型
// @Param dictName query string false "dictName"
// @Param dictId query string false "dictId"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/type/list [get]
// @Security
func GetDictTypeList(c *gin.Context) {
	var data models.DictType
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.DictName = c.Request.FormValue("dictName")
	id := c.Request.FormValue("dictId")
	data.DictId, _ = utils.StringToInt(id)
	data.DictType = c.Request.FormValue("dictType")
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

// @Summary 通过字典id获取字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Param dictId path int true "字典类型编码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/type/{dictId} [get]
// @Security
func GetDictType(c *gin.Context) {
	var DictType models.DictType
	DictType.DictName = c.Request.FormValue("dictName")
	DictType.DictId, _ = utils.StringToInt(c.Param("dictId"))
	result, err := DictType.Get()
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}


func GetDictTypeOptionSelect(c *gin.Context) {
	var DictType models.DictType
	DictType.DictName = c.Request.FormValue("dictName")
	DictType.DictId, _ = utils.StringToInt(c.Param("dictId"))
	result, err := DictType.GetList()
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 添加字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/type [post]
// @Security Bearer
func InsertDictType(c *gin.Context) {
	var data models.DictType
	err := c.BindWith(&data, binding.JSON)
	data.CreateBy = utils.GetUserIdStr(c)
	pkg.HasError(err, "", 500)
	result, err := data.Create()
	pkg.HasError(err, "", -1)
	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 修改字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/type [put]
// @Security Bearer
func UpdateDictType(c *gin.Context) {
	var data models.DictType
	err := c.BindWith(&data, binding.JSON)
	data.UpdateBy = utils.GetUserIdStr(c)
	pkg.HasError(err, "", -1)
	result, err := data.Update(data.DictId)
	pkg.HasError(err, "", -1)
	var res app.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 删除字典类型
// @Description 删除数据
// @Tags 字典类型
// @Param dictId path int true "dictId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dict/type/{dictId} [delete]
func DeleteDictType(c *gin.Context) {
	var data models.DictType
	id, err := utils.StringToInt(c.Param("dictId"))
	data.UpdateBy = utils.GetUserIdStr(c)
	_, err = data.Delete(id)
	pkg.HasError(err, "修改失败", 500)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
