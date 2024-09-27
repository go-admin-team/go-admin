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

type OrderReconciliation struct {
	api.Api
}

// GetPage 获取对账订单表，用于记录订单的对账信息列表
// @Summary 获取对账订单表，用于记录订单的对账信息列表
// @Description 获取对账订单表，用于记录订单的对账信息列表
// @Tags 对账订单表，用于记录订单的对账信息
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.OrderReconciliation}} "{"code": 200, "data": [...]}"
// @Router /api/v1/order-reconciliation [get]
// @Security Bearer
func (e OrderReconciliation) GetPage(c *gin.Context) {
	req := dto.OrderReconciliationGetPageReq{}
	s := service.OrderReconciliation{}
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
	list := make([]models.OrderReconciliation, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取对账订单表，用于记录订单的对账信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取对账订单表，用于记录订单的对账信息
// @Summary 获取对账订单表，用于记录订单的对账信息
// @Description 获取对账订单表，用于记录订单的对账信息
// @Tags 对账订单表，用于记录订单的对账信息
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.OrderReconciliation} "{"code": 200, "data": [...]}"
// @Router /api/v1/order-reconciliation/{id} [get]
// @Security Bearer
func (e OrderReconciliation) Get(c *gin.Context) {
	req := dto.OrderReconciliationGetReq{}
	s := service.OrderReconciliation{}
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
	var object models.OrderReconciliation

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取对账订单表，用于记录订单的对账信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建对账订单表，用于记录订单的对账信息
// @Summary 创建对账订单表，用于记录订单的对账信息
// @Description 创建对账订单表，用于记录订单的对账信息
// @Tags 对账订单表，用于记录订单的对账信息
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderReconciliationInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/order-reconciliation [post]
// @Security Bearer
func (e OrderReconciliation) Insert(c *gin.Context) {
	req := dto.OrderReconciliationInsertReq{}
	s := service.OrderReconciliation{}
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
		e.Error(500, err, fmt.Sprintf("创建对账订单表，用于记录订单的对账信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Merge 合并对账订单表，用于记录订单的对账信息(支付宝)
// @Summary 创建对账订单表，用于记录订单的对账信息（支付宝）
// @Description 创建对账订单表，用于记录订单的对账信息（支付宝）
// @Tags 对账订单表，用于记录订单的对账信息
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderReconciliationInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/order-reconciliation [post]
// @Security Bearer
func (e OrderReconciliation) Merge(c *gin.Context) {
	req := dto.OrderReconciliationInsertReq{}
	s := service.OrderReconciliation{}
	s.InitDB()
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
	req.Id = 1
	logNotify, err1 := s.MergeData()
	if err1 != nil {
		e.Error(500, err, fmt.Sprintf("(支付宝)合并对账订单表，用于记录订单的对账信息失败，\r\n失败信息 %s", err1.Error()))
		return
	}

	e.OK(logNotify, "创建成功")
}

// MergeWx 合并对账订单表，用于记录订单的对账信息(微信)
// @Summary 创建对账订单表，用于记录订单的对账信息（微信）
// @Description 创建对账订单表，用于记录订单的对账信息（微信）
// @Tags 对账订单表，用于记录订单的对账信息
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderReconciliationInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/order-reconciliation [post]
// @Security Bearer
func (e OrderReconciliation) MergeWx(c *gin.Context) {
	req := dto.OrderReconciliationInsertReq{}
	s := service.OrderReconciliation{}
	//s.InitDB()
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
	err1 := s.MergeDataWx(c)
	if err1 != nil {
		e.Error(500, err1, fmt.Sprintf("(微信)合并对账订单表，用于记录订单的对账信息失败，\r\n失败信息 %s", err1.Error()))
		return
	}
	e.OK("", "创建成功")
}

// ScanWx 微信扫码(微信)
// @Summary 微信扫码（微信）
// @Description 创微信扫码（微信）
// @Tags 微信扫码
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderReconciliationInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/order-reconciliation [post]
// @Security Bearer
func (e OrderReconciliation) ScanWx(c *gin.Context) {
	req := dto.OrderReconciliationInsertReq{}
	s := service.OrderReconciliation{}
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
	err1 := s.MergeDataWx(c)
	if err1 != nil {
		e.Error(500, err1, fmt.Sprintf("获取微信扫码信息失败，\r\n失败信息 %s", err1.Error()))
		return
	}
	e.OK("", "创建成功")
}

// Update 修改对账订单表，用于记录订单的对账信息
// @Summary 修改对账订单表，用于记录订单的对账信息
// @Description 修改对账订单表，用于记录订单的对账信息
// @Tags 对账订单表，用于记录订单的对账信息
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderReconciliationUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/order-reconciliation/{id} [put]
// @Security Bearer
func (e OrderReconciliation) Update(c *gin.Context) {
	req := dto.OrderReconciliationUpdateReq{}
	s := service.OrderReconciliation{}
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
		e.Error(500, err, fmt.Sprintf("修改对账订单表，用于记录订单的对账信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除对账订单表，用于记录订单的对账信息
// @Summary 删除对账订单表，用于记录订单的对账信息
// @Description 删除对账订单表，用于记录订单的对账信息
// @Tags 对账订单表，用于记录订单的对账信息
// @Param data body dto.OrderReconciliationDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/order-reconciliation [delete]
// @Security Bearer
func (e OrderReconciliation) Delete(c *gin.Context) {
	s := service.OrderReconciliation{}
	req := dto.OrderReconciliationDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除对账订单表，用于记录订单的对账信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
