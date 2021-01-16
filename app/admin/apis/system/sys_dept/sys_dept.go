package sys_dept

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/models"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
	"go-admin/common/log"
	"go-admin/tools"

	"net/http"
)

type SysDept struct {
	apis.Api
}

func (e *SysDept) GetSysDeptList(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	d := new(dto.SysDeptSearch)
	db, err := tools.GetOrm(c)
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

	list := make([]models.SysDept, 0)
	var count int64
	serviceStudent := service.SysDept{}
	serviceStudent.MsgID = msgID
	serviceStudent.Orm = db
	err = serviceStudent.GetSysDeptPage(d, &list, &count)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(c, list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e *SysDept) GetSysDept(c *gin.Context) {
	control := new(dto.SysDeptById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object models.SysDept

	serviceSysOperlog := service.SysDept{}
	serviceSysOperlog.MsgID = msgID
	serviceSysOperlog.Orm = db
	err = serviceSysOperlog.GetSysDept(control, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

func (e *SysDept) InsertSysDept(c *gin.Context) {
	control := new(dto.SysDeptControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
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
	object.SetCreateBy(tools.GetUserId(c))

	serviceSysDept := service.SysDept{}
	serviceSysDept.Orm = db
	serviceSysDept.MsgID = msgID
	err = serviceSysDept.InsertSysDept(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

func (e *SysDept) UpdateSysDept(c *gin.Context) {
	control := new(dto.SysDeptControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
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
	object.SetUpdateBy(tools.GetUserId(c))

	serviceSysDept := service.SysDept{}
	serviceSysDept.Orm = db
	serviceSysDept.MsgID = msgID
	err = serviceSysDept.UpdateSysDept(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

func (e *SysDept) DeleteSysDept(c *gin.Context) {
	control := new(dto.SysDeptById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Errorf("MsgID[%s] Bind error: %s", msgID, err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysDept := service.SysDept{}
	serviceSysDept.Orm = db
	serviceSysDept.MsgID = msgID
	err = serviceSysDept.RemoveSysDept(control)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.GetId(), "删除成功")
}
