package apis

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysPost struct {
	api.Api
}

// GetPage
// @Summary 岗位列表数据
// @Description 获取JSON
// @Tags 岗位
// @Param postName query string false "postName"
// @Param postCode query string false "postCode"
// @Param postId query string false "postId"
// @Param status query string false "status"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/post [get]
// @Security Bearer
func (e SysPost) GetPage(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostSearch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	list := make([]models.SysPost, 0)
	var count int64

	err = s.GetSysPostPage(&req, &list, &count)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get
// @Summary 获取岗位信息
// @Description 获取JSON
// @Tags 岗位
// @Param postId path int true "postId"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/post/{postId} [get]
// @Security Bearer
func (e SysPost) Get(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.SysPost

	err = s.GetSysPost(&req, &object)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.OK(object, "查询成功")
}

// Insert
// @Summary 添加岗位
// @Description 获取JSON
// @Tags 岗位
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysPostControl true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/post [post]
// @Security Bearer
func (e SysPost) Insert(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetCreateBy(user.GetUserId(c))
	err = s.InsertSysPost(&req)
	if err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.GetId(), "创建成功")
}

// Update
// @Summary 修改岗位
// @Description 获取JSON
// @Tags 岗位
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysPostControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/post/ [put]
// @Security Bearer
func (e SysPost) Update(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	req.SetUpdateBy(user.GetUserId(c))

	err = s.UpdateSysPost(&req)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// Delete
// @Summary 删除岗位
// @Description 删除数据
// @Tags 岗位
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 500 {string} string	"{"code": 500, "message": "删除失败"}"
// @Router /api/v1/post/{postId} [delete]
func (e SysPost) Delete(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	err = s.RemoveSysPost(&req)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "删除成功")
}