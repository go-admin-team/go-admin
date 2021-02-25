package syschinaareadata

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/service"
	"net/http"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
	"go-admin/common/log"
	"go-admin/tools"
)

type SysChinaAreaData struct {
	apis.Api
}

func (e *SysChinaAreaData) GetSysChinaAreaDataList(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	d := new(dto.SysChinaAreaDataSearch)
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

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysChinaAreaData, 0)
	var count int64
	serviceStudent := service.SysChinaAreaData{}
	serviceStudent.MsgID = msgID
	serviceStudent.Orm = db
	err = serviceStudent.GetSysChinaAreaDataPage(d, p, &list, &count)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(c, list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e *SysChinaAreaData) GetSysChinaAreaData(c *gin.Context) {
	control := new(dto.SysChinaAreaDataById)
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
	var object models.SysChinaAreaData

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.MsgID = msgID
	serviceSysChinaAreaData.Orm = db
	err = serviceSysChinaAreaData.GetSysChinaAreaData(control, p, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

func (e *SysChinaAreaData) InsertSysChinaAreaData(c *gin.Context) {
	control := new(dto.SysChinaAreaDataControl)
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

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.Orm = db
	serviceSysChinaAreaData.MsgID = msgID
	err = serviceSysChinaAreaData.InsertSysChinaAreaData(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

func (e *SysChinaAreaData) UpdateSysChinaAreaData(c *gin.Context) {
	control := new(dto.SysChinaAreaDataControl)
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

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.Orm = db
	serviceSysChinaAreaData.MsgID = msgID
	err = serviceSysChinaAreaData.UpdateSysChinaAreaData(object, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

func (e *SysChinaAreaData) DeleteSysChinaAreaData(c *gin.Context) {
	control := new(dto.SysChinaAreaDataById)
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

	// 设置编辑人
	control.SetUpdateBy(tools.GetUserId(c))

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.Orm = db
	serviceSysChinaAreaData.MsgID = msgID
	err = serviceSysChinaAreaData.RemoveSysChinaAreaData(control, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.GetId(), "删除成功")
}
