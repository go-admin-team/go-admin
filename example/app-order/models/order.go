// Package models holds app-order's two GORM row models.
package models

import (
	contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"
)

// The two values Order.Status can hold. Kept as narrow strings rather than
// an int enum to match sys_role.data_scope's own convention in core, and to
// leave room for a future status without a schema change.
const (
	StatusPending = "1" // awaiting payment
	StatusPaid    = "2" // paid; set only by a successful Pay
)

// orderTable is passed to actions.Permission and repeated as
// Order.TableName's return value. It is not literally the word "order":
// that is a reserved SQL keyword, and actions.Permission builds its WHERE
// clause by string-concatenating tableName straight into raw SQL
// (`tableName+".create_by = ?"`, see permission.go) with no quoting at all.
// A table named exactly "order" would make every data-scope query a syntax
// error on MySQL's default (non-ANSI-quotes) mode. This is not something
// core enforces or even mentions - Permission's tableName parameter is an
// opaque string as far as it is concerned - so avoiding reserved words is
// entirely on the caller.
const orderTable = "app_order"

// Order is one customer order. ControlBy is required, not decorative:
// actions.Permission's data-scope SQL joins against create_by, so an Order
// without it would make every data-scope rule silently match nothing.
type Order struct {
	contractmodels.Model

	OrderNo    string `json:"orderNo" gorm:"type:varchar(64);uniqueIndex;comment:order number"`
	UserId     int    `json:"userId" gorm:"index;comment:buyer user id"`
	Status     string `json:"status" gorm:"type:varchar(4);index;comment:order status: 1 pending, 2 paid"`
	TotalCents int64  `json:"totalCents" gorm:"comment:total amount in cents, sum of item price*quantity at creation time"`

	// Items is populated by Preload; it is never set by Order's own migrator
	// column set (OrderItem.OrderId is the foreign key, not a column here).
	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderId"`

	contractmodels.ControlBy
	contractmodels.ModelTime
}

// TableName pins the row model to app_order regardless of any global
// singular/plural table naming strategy the host configures. See orderTable
// above for why this is not simply "order".
func (Order) TableName() string {
	return orderTable
}
