package tools

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/tools"
	"go-admin/pkg"
	"net/http"
)

// @Summary 分页列表数据 / page list data
// @Description 生成表分页列表 / gen table page list
// @Tags 工具 - Tools / 生成表 - Gen Table
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys/tables/page [get]
func GetSysTableList(c *gin.Context) {
	var data tools.SysTables
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.TableName = c.Request.FormValue("tableName")
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.AssertErr(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageIndex"] = pageSize

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}