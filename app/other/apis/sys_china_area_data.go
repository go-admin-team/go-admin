package apis

import (
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/other/models"
	"go-admin/app/other/service"
	"go-admin/app/other/service/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/common/actions"
)

type SysChinaAreaData struct {
	api.Api
}

func (e SysChinaAreaData) GetPage(c *gin.Context) {
	s := service.SysChinaAreaData{}
	req := dto.SysChinaAreaDataSearch{}
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

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysChinaAreaData, 0)
	var count int64
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e SysChinaAreaData) Get(c *gin.Context) {
	s := service.SysChinaAreaData{}
	req := dto.SysChinaAreaDataById{}
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
	var object models.SysChinaAreaData
	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查询成功")
}

func (e SysChinaAreaData) Insert(c *gin.Context) {
	e.MakeContext(c)
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
	err = serviceSysChinaAreaData.Insert(object)
	if err != nil {
		log.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(object.GetId(), "创建成功")
}

func (e SysChinaAreaData) Update(c *gin.Context) {
	e.MakeContext(c)
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
	err = serviceSysChinaAreaData.Update(object, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(object.GetId(), "更新成功")
}

func (e SysChinaAreaData) Delete(c *gin.Context) {
	e.MakeContext(c)
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
	err = serviceSysChinaAreaData.Remove(control, p)
	if err != nil {
		log.Errorf("RemoveSysChinaAreaData error, %s", err)
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(control.GetId(), "删除成功")
}
