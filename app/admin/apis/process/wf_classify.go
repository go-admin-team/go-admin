package process

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
)

type WfProcessClassify struct {
	apis.Api
}

func (e *WfProcessClassify) GetWfProcessClassifyList(c *gin.Context) {
	log := e.GetLogger(c)
	d := new(dto.WfProcessClassifySearch)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查询列表
	err = d.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.WfProcessClassify, 0)
	var count int64
	serviceStudent := service.WfProcessClassify{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetWfProcessClassifyPage(d, p, &list, &count)
	if err != nil {
		log.Errorf("GetWfProcessClassifyPage error, %s", err)
		e.Error(c, http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.PageOK(c, list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e *WfProcessClassify) GetWfProcessClassify(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.WfProcessClassifyById)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object models.WfProcessClassify

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceWfProcessClassify := service.WfProcessClassify{}
	serviceWfProcessClassify.Log = log
	serviceWfProcessClassify.Orm = db
	err = serviceWfProcessClassify.GetWfProcessClassify(control, p, &object)
	if err != nil {
		log.Errorf("GetWfProcessClassify error, %s", err)
		e.Error(c, http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

func (e *WfProcessClassify) InsertWfProcessClassify(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.WfProcessClassifyControl)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//新增操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	// 设置创建人
	object.SetCreateBy(user.GetUserId(c))

	serviceWfProcessClassify := service.WfProcessClassify{}
	serviceWfProcessClassify.Orm = db
	serviceWfProcessClassify.Log = log
	err = serviceWfProcessClassify.InsertWfProcessClassify(object)
	if err != nil {
		log.Errorf("InsertWfProcessClassify error, %s", err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

func (e *WfProcessClassify) UpdateWfProcessClassify(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.WfProcessClassifyControl)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//更新操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.SetUpdateBy(user.GetUserId(c))

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceWfProcessClassify := service.WfProcessClassify{}
	serviceWfProcessClassify.Orm = db
	serviceWfProcessClassify.Log = log
	err = serviceWfProcessClassify.UpdateWfProcessClassify(object, p)
	if err != nil {
		log.Errorf("UpdateWfProcessClassify error, %s", err)
		e.Error(c, http.StatusInternalServerError, err, "更新失败")
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

func (e *WfProcessClassify) DeleteWfProcessClassify(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.WfProcessClassifyById)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Warnf("Bind error: %s", err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	// 设置编辑人
	control.SetUpdateBy(user.GetUserId(c))

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceWfProcessClassify := service.WfProcessClassify{}
	serviceWfProcessClassify.Orm = db
	serviceWfProcessClassify.Log = log
	err = serviceWfProcessClassify.RemoveWfProcessClassify(control, p)
	if err != nil {
		log.Errorf("RemoveWfProcessClassify error, %s", err)
		e.Error(c, http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(c, control.GetId(), "删除成功")
}
