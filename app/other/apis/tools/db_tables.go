package tools

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/other/models/tools"
)

// GetDBTableList 分页列表数据
// @Summary 分页列表数据 / page list data
// @Description 数据库表分页列表 / database table page list
// @Tags 工具 / 生成工具
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/db/tables/page [get]
func (e *Gen) GetDBTableList(c *gin.Context) {
	//var res response.Response
	var data tools.DBTables
	var err error
	var pageSize = 10
	var pageIndex = 1
	e.Context = c
	log := e.GetLogger()
	if config.DatabaseConfig.Driver == "sqlite3" || config.DatabaseConfig.Driver == "postgres" {
		err = errors.New("对不起，sqlite3 或 postgres 不支持代码生成！")
		log.Warn(err)
		e.Error(403, err, "")
		return
	}

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = pkg.StringToInt(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = pkg.StringToInt(index)
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(500, err, "数据库连接获取失败")
		return
	}

	data.TableName = c.Request.FormValue("tableName")
	result, count, err := data.GetPage(db, pageSize, pageIndex)
	if err != nil {
		log.Errorf("GetPage error, %s", err.Error())
		e.Error(500, err, "")
		return
	}
	e.PageOK(result, count, pageIndex, pageSize, "查询成功")
}
