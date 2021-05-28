package sys_post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
)

type SysPost struct {
	apis.Api
}

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
func (e SysPost) GetSysPostList(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	d := new(dto.SysPostSearch)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//查询列表
	err = d.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	list := make([]system.SysPost, 0)
	var count int64
	serviceStudent := service.SysPost{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetSysPostPage(d, &list, &count)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

// @Summary 获取岗位信息
// @Description 获取JSON
// @Tags 岗位
// @Param postId path int true "postId"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/post/{postId} [get]
// @Security Bearer
func (e SysPost) GetSysPost(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysPostById)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object system.SysPost

	serviceSysOperlog := service.SysPost{}
	serviceSysOperlog.Log = log
	serviceSysOperlog.Orm = db
	err = serviceSysOperlog.GetSysPost(control, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

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
func (e SysPost) InsertSysPost(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysPostControl)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//新增操作
	err = control.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	// 设置创建人
	object.SetCreateBy(user.GetUserId(c))

	serviceSysPost := service.SysPost{}
	serviceSysPost.Orm = db
	serviceSysPost.Log = log
	err = serviceSysPost.InsertSysPost(object)
	if err != nil {
		log.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(object.GetId(), "创建成功")
}

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
func (e SysPost) UpdateSysPost(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysPostControl)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//更新操作
	err = control.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.SetUpdateBy(user.GetUserId(c))

	serviceSysPost := service.SysPost{}
	serviceSysPost.Orm = db
	serviceSysPost.Log = log
	err = serviceSysPost.UpdateSysPost(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(object.GetId(), "更新成功")
}

// @Summary 删除岗位
// @Description 删除数据
// @Tags 岗位
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 500 {string} string	"{"code": 500, "message": "删除失败"}"
// @Router /api/v1/post/{postId} [delete]
func (e SysPost) DeleteSysPost(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysPostById)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Errorf("Bind error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysPost := service.SysPost{}
	serviceSysPost.Orm = db
	serviceSysPost.Log = log
	err = serviceSysPost.RemoveSysPost(control)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(control.GetId(), "删除成功")
}
