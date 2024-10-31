package apis

import (
    "fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

type TApiZl struct {
	api.Api
}

// GetPage 获取TApiZl列表
// @Summary 获取TApiZl列表
// @Description 获取TApiZl列表
// @Tags TApiZl
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.TApiZl}} "{"code": 200, "data": [...]}"
// @Router /api/v1/t-api-zl [get]
// @Security Bearer
func (e TApiZl) GetPage(c *gin.Context) {
    req := dto.TApiZlGetPageReq{}
    s := service.TApiZl{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
   	if err != nil {
   		e.Logger.Error(err)
   		e.Error(500, err, err.Error())
   		return
   	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.TApiZl, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取TApiZl失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取TApiZl
// @Summary 获取TApiZl
// @Description 获取TApiZl
// @Tags TApiZl
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.TApiZl} "{"code": 200, "data": [...]}"
// @Router /api/v1/t-api-zl/{id} [get]
// @Security Bearer
func (e TApiZl) Get(c *gin.Context) {
	req := dto.TApiZlGetReq{}
	s := service.TApiZl{}
    err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.TApiZl

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取TApiZl失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建TApiZl
// @Summary 创建TApiZl
// @Description 创建TApiZl
// @Tags TApiZl
// @Accept application/json
// @Product application/json
// @Param data body dto.TApiZlInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/t-api-zl [post]
// @Security Bearer
func (e TApiZl) Insert(c *gin.Context) {
    req := dto.TApiZlInsertReq{}
    s := service.TApiZl{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建TApiZl失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改TApiZl
// @Summary 修改TApiZl
// @Description 修改TApiZl
// @Tags TApiZl
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.TApiZlUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/t-api-zl/{id} [put]
// @Security Bearer
func (e TApiZl) Update(c *gin.Context) {
    req := dto.TApiZlUpdateReq{}
    s := service.TApiZl{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改TApiZl失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除TApiZl
// @Summary 删除TApiZl
// @Description 删除TApiZl
// @Tags TApiZl
// @Param data body dto.TApiZlDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/t-api-zl [delete]
// @Security Bearer
func (e TApiZl) Delete(c *gin.Context) {
    s := service.TApiZl{}
    req := dto.TApiZlDeleteReq{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除TApiZl失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
