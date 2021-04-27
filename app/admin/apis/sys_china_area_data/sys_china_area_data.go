package sys_china_area_data

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

type SysChinaAreaData struct {
	apis.Api
}

func (e SysChinaAreaData) GetSysChinaAreaDataList(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	d := new(dto.SysChinaAreaDataSearch)
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

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysChinaAreaData, 0)
	var count int64
	serviceStudent := service.SysChinaAreaData{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetSysChinaAreaDataPage(d, p, &list, &count)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e SysChinaAreaData) GetSysChinaAreaData(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysChinaAreaDataById)
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
	var object models.SysChinaAreaData

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.Log = log
	serviceSysChinaAreaData.Orm = db
	err = serviceSysChinaAreaData.GetSysChinaAreaData(control, p, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

func (e SysChinaAreaData) InsertSysChinaAreaData(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysChinaAreaDataControl)
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

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.Orm = db
	serviceSysChinaAreaData.Log = log
	err = serviceSysChinaAreaData.InsertSysChinaAreaData(object)
	if err != nil {
		log.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(object.GetId(), "创建成功")
}

func (e SysChinaAreaData) UpdateSysChinaAreaData(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysChinaAreaDataControl)
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

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.Orm = db
	serviceSysChinaAreaData.Log = log
	err = serviceSysChinaAreaData.UpdateSysChinaAreaData(object, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(object.GetId(), "更新成功")
}

func (e SysChinaAreaData) DeleteSysChinaAreaData(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysChinaAreaDataById)
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

	// 设置编辑人
	control.SetUpdateBy(user.GetUserId(c))

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysChinaAreaData := service.SysChinaAreaData{}
	serviceSysChinaAreaData.Orm = db
	serviceSysChinaAreaData.Log = log
	err = serviceSysChinaAreaData.RemoveSysChinaAreaData(control, p)
	if err != nil {
		log.Errorf("RemoveSysChinaAreaData error, %s", err)
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(control.GetId(), "删除成功")
}
