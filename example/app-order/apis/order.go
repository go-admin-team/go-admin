// Package apis is app-order's HTTP layer: four hand-written gin handlers,
// none of them a wrapper around core's generic CRUD Actions. See
// service/order.go's package doc for why.
package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-admin-team/go-admin-core/v2/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/actions"

	// models.Response in the @Success annotations below resolves to
	// go-admin-core's sdk/contract/models.Response, not to this package -
	// swaggo finds it through --parseDependency. It is the envelope with a
	// data field; core's response.Response, which the framework's own
	// handlers name, has no data field and would document these endpoints
	// as returning none.
	"github.com/go-admin-team/example-app-order/models"
	"github.com/go-admin-team/example-app-order/service"
	orderdto "github.com/go-admin-team/example-app-order/service/dto"
)

// Order embeds api.Api the same way every hand-written go-admin handler
// does (see app/admin/apis/sys_post.go): MakeContext/MakeOrm/Bind/
// MakeService/OK/Error/PageOK are all core, imported with no dependency on
// go-admin itself.
type Order struct {
	api.Api
}

// GetPage
// @Summary List orders visible to the caller's data scope
// @Tags order
// @Param status query string false "status"
// @Param orderNo query string false "orderNo"
// @Param pageIndex query int false "pageIndex"
// @Param pageSize query int false "pageSize"
// @Success 200 {object} models.Response
// @Router /api/v1/order [get]
// @Security Bearer
func (e Order) GetPage(c *gin.Context) {
	s := service.Order{}
	req := orderdto.OrderSearchReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.Order, 0)
	count, err := s.GetPage(&req, p, &list)
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, "failed to list orders")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "ok")
}

// Get
// @Summary Get one order and its items
// @Tags order
// @Param id path int true "order id"
// @Success 200 {object} models.Response
// @Router /api/v1/order/{id} [get]
// @Security Bearer
func (e Order) Get(c *gin.Context) {
	s := service.Order{}
	req := orderdto.OrderIdReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	var order models.Order
	if err = s.Get(req.Id, p, &order); err != nil {
		e.Error(http.StatusNotFound, err, "order not found")
		return
	}
	e.OK(order, "ok")
}

// Create
// @Summary Place a new order
// @Tags order
// @Accept application/json
// @Param data body orderdto.OrderCreateReq true "data"
// @Success 200 {object} models.Response
// @Router /api/v1/order [post]
// @Security Bearer
func (e Order) Create(c *gin.Context) {
	s := service.Order{}
	req := orderdto.OrderCreateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, err.Error())
		return
	}

	order, err := s.Create(&req, user.GetUserId(c))
	if err != nil {
		if errors.Is(err, service.ErrOrderEmpty) {
			e.Error(http.StatusBadRequest, err, err.Error())
			return
		}
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, "failed to create order")
		return
	}
	e.OK(order, "created")
}

// Pay
// @Summary Mark a pending order as paid
// @Tags order
// @Param id path int true "order id"
// @Success 200 {object} models.Response
// @Router /api/v1/order/{id}/pay [put]
// @Security Bearer
func (e Order) Pay(c *gin.Context) {
	s := service.Order{}
	req := orderdto.OrderIdReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	if err = s.Pay(req.Id, p); err != nil {
		if errors.Is(err, service.ErrOrderNotPending) {
			// Deliberately the same response whether the order does not
			// exist, is already paid, or is outside p's data scope - see
			// service.Order.Pay's doc comment.
			e.Error(http.StatusConflict, err, err.Error())
			return
		}
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, "payment failed")
		return
	}
	e.OK(nil, "paid")
}
