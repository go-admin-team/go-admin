// Package dto holds app-order's request-binding types.
//
// None of them implement core's dto.Index / dto.Control, and none of them
// define their own Bind method: those interfaces (and the Bind method they
// require) exist so the framework's generic CRUD Actions
// (Create/Delete/Index/Update/ViewAction) can bind a request without
// knowing its concrete type - Action itself calls req.Bind(c). app-order's
// handlers (apis/order.go) call api.Api.Bind directly on the raw struct
// instead, exactly as go-admin's own hand-written handlers do (see
// app/admin/apis/sys_post.go and its service/dto/sys_post.go, which is the
// same shape: plain structs, no Bind method), so a Bind method here would
// never be called by anything and would only mislead a reader into thinking
// it is.
//
// What these types do reuse is dto.Pagination (for the list request's page
// index/size) and the `search` struct-tag convention dto.MakeCondition
// resolves; both are plain data shapes, not an interface a hand-written
// handler would otherwise have to reimplement.
package dto

import (
	contractdto "github.com/go-admin-team/go-admin-core/v2/sdk/contract/dto"
)

// OrderItemReq is one line item in a create-order request.
type OrderItemReq struct {
	ProductName string `json:"productName" validate:"required"`
	Quantity    int    `json:"quantity" validate:"gte=1"`
	PriceCents  int64  `json:"priceCents" validate:"gte=0"`
}

// OrderCreateReq is the create-order request body.
type OrderCreateReq struct {
	Items []OrderItemReq `json:"items" validate:"required"`
}

// OrderSearchReq is the list-order query.
//
// contractdto.MakeCondition reads q's `search` tags through
// reflect.TypeOf(q).NumField(), which is only valid for a struct Kind - a
// pointer panics rather than returning an error (see
// service/order.go:GetPage, which is careful to pass *req, not req). That
// distinction is not documented on MakeCondition's exported doc comment.
// Pagination `search:"-"` here follows the same convention the framework's
// own generic DTOs use to keep Pagination's two fields out of the WHERE
// clause the tags on Status/OrderNo build.
type OrderSearchReq struct {
	contractdto.Pagination `search:"-"`

	Status  string `form:"status" search:"type:exact;column:status;table:app_order"`
	OrderNo string `form:"orderNo" search:"type:exact;column:order_no;table:app_order"`
}

// OrderIdReq binds a single :id, for a detail lookup or a Pay request.
type OrderIdReq struct {
	Id int `uri:"id" validate:"required"`
}
