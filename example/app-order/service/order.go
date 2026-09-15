// Package service is app-order's business logic: everything the PRD asked
// this example to prove out by hand rather than by wiring up core's generic
// CRUD Actions (Create/Delete/Index/Update/ViewAction stay in go-admin, not
// in core - see contract/actions's package doc for why). An order's write
// path is a cross-table transaction and its one state change needs a
// concurrency guard neither generic Action was ever built for, which is
// exactly the class of logic real third-party apps almost always have.
package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/actions"
	contractdto "github.com/go-admin-team/go-admin-core/v2/sdk/contract/dto"
	"github.com/go-admin-team/go-admin-core/v2/sdk/service"

	"github.com/go-admin-team/example-app-order/models"
	orderdto "github.com/go-admin-team/example-app-order/service/dto"
)

// ErrOrderEmpty is returned by Create when the request has no line items.
var ErrOrderEmpty = errors.New("app-order: an order must have at least one item")

// ErrOrderNotPending is returned by Pay when the order could not be paid:
// it does not exist, it is not in models.StatusPending, or the caller's
// data scope does not include it. Deliberately one error for all three -
// see Pay's doc comment for why collapsing them is the fail-closed choice,
// not a shortcut.
var ErrOrderNotPending = errors.New("app-order: order is not awaiting payment")

// Order is app-order's hand-written service. It embeds core's
// sdk/service.Service purely for the Orm/Log/Cache fields every
// api.Api.MakeService caller already wires up the same way go-admin's own
// hand-written services do (see app/admin/apis/sys_post.go) - not because
// anything here calls a method Service defines.
type Order struct {
	service.Service
}

// Create places a new order. The order row and every item row commit
// together: db.Transaction's closure form is what makes that true even
// across a panic (it recovers, rolls back, and re-panics - see gorm's own
// Transaction implementation), unlike the hand-rolled Begin/defer pattern
// go-admin's sys_role.go/sys_dept.go/sys_menu.go/sys_tables.go use, which
// commits a half-written transaction on panic, never opens a real
// transaction under sqlite, and reads a single global DB handle regardless
// of which tenant the request is for.
func (e *Order) Create(req *orderdto.OrderCreateReq, userId int) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, ErrOrderEmpty
	}

	items := make([]models.OrderItem, 0, len(req.Items))
	var total int64
	for _, it := range req.Items {
		total += it.PriceCents * int64(it.Quantity)
		items = append(items, models.OrderItem{
			ProductName: it.ProductName,
			Quantity:    it.Quantity,
			PriceCents:  it.PriceCents,
		})
	}

	order := &models.Order{
		OrderNo:    generateOrderNo(),
		UserId:     userId,
		Status:     models.StatusPending,
		TotalCents: total,
	}
	order.SetCreateBy(userId)
	order.SetUpdateBy(userId)

	err := e.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderId = order.Id
		}
		// A single batch Create, not one Create per item: on the unique
		// (order_id, product_name) violation the test suite exercises, the
		// whole statement fails, and nothing about this order - not the
		// order row created two lines above, not any item before the
		// duplicate - survives the rollback.
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	order.Items = items
	return order, nil
}

// Get loads one order, scoped to p's data permission, with its items.
func (e *Order) Get(id int, p *actions.DataPermission, out *models.Order) error {
	return e.Orm.
		Scopes(actions.Permission(orderTableName, p)).
		Preload("Items").
		Where("id = ?", id).
		First(out).Error
}

// GetPage lists orders visible to p's data scope, filtered by req's search
// tags and paginated. The Find-then-Count-on-the-same-chain shape mirrors
// go-admin's own common/actions.IndexAction: Limit(-1).Offset(-1) undoes
// Paginate's LIMIT/OFFSET before the count runs, on the same *gorm.DB
// session, so the WHERE clause built by MakeCondition and Permission is not
// re-resolved a second time.
func (e *Order) GetPage(req *orderdto.OrderSearchReq, p *actions.DataPermission, list *[]models.Order) (int64, error) {
	var count int64
	// *req, not req: contractdto.MakeCondition resolves search tags through
	// reflect.TypeOf(q).NumField(), which panics on a pointer. See
	// service/dto/order.go's doc comment on OrderSearchReq.
	err := e.Orm.Model(&models.Order{}).
		Scopes(
			contractdto.MakeCondition(*req),
			contractdto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(orderTableName, p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(&count).Error
	return count, err
}

// Pay transitions a pending order to paid.
//
// The concurrency guard is the WHERE clause, not an application-level lock:
// two concurrent payment attempts against the same order both issue this
// UPDATE, but only the one that actually flips a row from pending to paid
// sees RowsAffected == 1 - the loser's WHERE matches nothing (the row is
// already 'paid' by the time its UPDATE runs) and sees 0, becoming
// ErrOrderNotPending rather than a second, silently-accepted payment.
//
// The same RowsAffected==0 outcome also covers "no such order" and "this
// order exists but is outside p's data scope" - actions.Permission's own
// scope is one of the Scopes below, so a caller paying an order they
// cannot see gets the identical error a caller paying an already-paid
// order gets. That collapse is deliberate: a distinguishable "exists but
// not yours" response would leak which order ids exist to a caller who
// should not be able to tell.
func (e *Order) Pay(id int, p *actions.DataPermission) error {
	result := e.Orm.
		Scopes(actions.Permission(orderTableName, p)).
		Model(&models.Order{}).
		Where("id = ? AND status = ?", id, models.StatusPending).
		Updates(map[string]interface{}{"status": models.StatusPaid})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOrderNotPending
	}
	return nil
}

// orderTableName is models.Order{}.TableName(), repeated here as a plain
// string because actions.Permission takes the table name as a bare string,
// not a model - see models/order.go's orderTable doc comment for why it is
// not literally "order".
const orderTableName = "app_order"

// generateOrderNo is a placeholder good enough for this example: real
// production code would want a collision-proof id source (a sequence, a
// snowflake id, or similar). Nothing about the transaction or the
// concurrency guard above depends on how this string is built.
func generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}
