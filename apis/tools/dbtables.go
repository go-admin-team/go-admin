package tools

import (
	"github.com/gin-gonic/gin"
	"go-admin/models/tools"
	tools2 "go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/config"
	"net/http"
)

// @Summary 分页列表数据 / page list data
// @Description 数据库表分页列表 / database table page list
// @Tags 工具 / Tools
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/db/tables/page [get]
func GetDBTableList(c *gin.Context) {
	var res app.Response
	var data tools.DBTables
	var err error
	var pageSize = 10
	var pageIndex = 1
	if config.DatabaseConfig.Driver == "sqlite3" || config.DatabaseConfig.Driver == "postgres" {
		res.Msg = "对不起，sqlite3 或 postgres 不支持代码生成！"
		c.JSON(http.StatusOK, res.ReturnError(500))
		return
	}

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools2.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools2.StrToInt(err, index)
	}

	data.TableName = c.Request.FormValue("tableName")
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools2.HasError(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}
