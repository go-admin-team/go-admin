package service

import (
	"errors"
	"sync"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/actions"
	coreservice "github.com/go-admin-team/go-admin-core/v2/sdk/service"

	"github.com/go-admin-team/example-app-order/models"
	orderdto "github.com/go-admin-team/example-app-order/service/dto"
)

// testDB returns a fresh, isolated in-memory sqlite database with
// app_order/app_order_item created, following the same
// glebarez/sqlite-and-no-build-tag setup core's own contract package tests
// use (see sdk/contract/actions/permission_test.go and
// sdk/contract/seed/seed_test.go).
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Order{}, &models.OrderItem{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

// enableDataPermission flips on the switch actions.Permission checks before
// applying any data-scope filtering at all, restoring the previous value
// after the test - the same pattern
// sdk/contract/actions/permission_test.go uses.
func enableDataPermission(t *testing.T) {
	t.Helper()
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = previous })
}

func newOrderService(t *testing.T, db *gorm.DB) *Order {
	t.Helper()
	return &Order{Service: coreservice.Service{Orm: db}}
}

// -- cross-table transaction ------------------------------------------------

func TestCreate_CommitsOrderAndItemsTogether(t *testing.T) {
	db := testDB(t)
	s := newOrderService(t, db)

	req := &orderdto.OrderCreateReq{Items: []orderdto.OrderItemReq{
		{ProductName: "widget", Quantity: 2, PriceCents: 500},
		{ProductName: "gadget", Quantity: 1, PriceCents: 1200},
	}}

	order, err := s.Create(req, 42)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if order.TotalCents != 2*500+1200 {
		t.Errorf("TotalCents = %d, want %d", order.TotalCents, 2*500+1200)
	}
	if order.Status != models.StatusPending {
		t.Errorf("Status = %q, want pending", order.Status)
	}
	if order.CreateBy != 42 || order.UpdateBy != 42 {
		t.Errorf("CreateBy/UpdateBy = %d/%d, want 42/42", order.CreateBy, order.UpdateBy)
	}

	var itemCount int64
	db.Model(&models.OrderItem{}).Where("order_id = ?", order.Id).Count(&itemCount)
	if itemCount != 2 {
		t.Errorf("persisted %d items, want 2", itemCount)
	}
}

func TestCreate_EmptyItemsReturnsErrorAndWritesNothing(t *testing.T) {
	db := testDB(t)
	s := newOrderService(t, db)

	_, err := s.Create(&orderdto.OrderCreateReq{}, 1)
	if !errors.Is(err, ErrOrderEmpty) {
		t.Fatalf("got error %v, want ErrOrderEmpty", err)
	}

	var count int64
	db.Model(&models.Order{}).Count(&count)
	if count != 0 {
		t.Errorf("an order was written despite the empty-items error")
	}
}

// A mid-transaction failure must roll back everything written before it in
// the same transaction, including the parent row. The duplicate product
// name is what forces a real, DB-enforced constraint violation on the
// second item's insert - see OrderItem's doc comment.
func TestCreate_MidTransactionFailureRollsBackEverything(t *testing.T) {
	db := testDB(t)
	s := newOrderService(t, db)

	req := &orderdto.OrderCreateReq{Items: []orderdto.OrderItemReq{
		{ProductName: "widget", Quantity: 1, PriceCents: 100},
		{ProductName: "widget", Quantity: 1, PriceCents: 100}, // duplicate: violates uk_app_order_item_product
	}}

	_, err := s.Create(req, 1)
	if err == nil {
		t.Fatal("Create succeeded despite a duplicate line item; the unique constraint did not fire")
	}

	var orderCount, itemCount int64
	db.Model(&models.Order{}).Count(&orderCount)
	db.Model(&models.OrderItem{}).Count(&itemCount)
	if orderCount != 0 {
		t.Errorf("the order row survived the rollback: %d rows in app_order", orderCount)
	}
	if itemCount != 0 {
		t.Errorf("an item row survived the rollback: %d rows in app_order_item", itemCount)
	}
}

// A panic partway through the transaction must roll back exactly as
// cleanly as a returned error does. This is not testing app-order's own
// code so much as the primitive Create is built on: gorm's db.Transaction
// recovers a panic, rolls back, and re-panics, which is what makes it safe
// to use in place of go-admin's hand-rolled Begin/defer pattern (see
// Create's doc comment) - a pattern that, on a panic, commits whatever the
// transaction had written so far instead of undoing it.
func TestCreate_PanicInsideTransactionRollsBackEverything(t *testing.T) {
	db := testDB(t)

	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("db.Transaction did not propagate the panic")
			}
		}()
		_ = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&models.Order{OrderNo: "panic-test", Status: models.StatusPending}).Error; err != nil {
				t.Fatalf("Create inside transaction: %v", err)
			}
			panic("simulated failure after a partial write")
		})
	}()

	var count int64
	db.Model(&models.Order{}).Count(&count)
	if count != 0 {
		t.Errorf("the order row survived a panic mid-transaction: %d rows in app_order", count)
	}
}

// -- status transition / concurrency guard ----------------------------------

func createPendingOrder(t *testing.T, s *Order, userId int) *models.Order {
	t.Helper()
	order, err := s.Create(&orderdto.OrderCreateReq{Items: []orderdto.OrderItemReq{
		{ProductName: "widget", Quantity: 1, PriceCents: 100},
	}}, userId)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	return order
}

func TestPay_TransitionsPendingToPaid(t *testing.T) {
	db := testDB(t)
	s := newOrderService(t, db)
	order := createPendingOrder(t, s, 1)

	if err := s.Pay(order.Id, &actions.DataPermission{DataScope: actions.DataScopeAll}); err != nil {
		t.Fatalf("Pay: %v", err)
	}

	var got models.Order
	db.First(&got, order.Id)
	if got.Status != models.StatusPaid {
		t.Errorf("Status = %q, want paid", got.Status)
	}
}

func TestPay_AlreadyPaidReturnsErrOrderNotPending(t *testing.T) {
	db := testDB(t)
	s := newOrderService(t, db)
	order := createPendingOrder(t, s, 1)
	all := &actions.DataPermission{DataScope: actions.DataScopeAll}

	if err := s.Pay(order.Id, all); err != nil {
		t.Fatalf("first Pay: %v", err)
	}
	if err := s.Pay(order.Id, all); !errors.Is(err, ErrOrderNotPending) {
		t.Fatalf("second Pay returned %v, want ErrOrderNotPending", err)
	}
}

// Two concurrent payment attempts against the same pending order: exactly
// one must succeed. MaxOpenConns(1) is set on the underlying *sql.DB so the
// two goroutines' UPDATEs serialize the way two independent connections
// would under MySQL, rather than one of them failing outright with
// SQLITE_BUSY - sqlite is a single-writer database with no useful
// concurrency of its own to exercise here. What the test actually verifies
// is unaffected by that: the guard is the UPDATE ... WHERE status =
// 'pending' clause and the RowsAffected check on its result (Pay's doc
// comment), and that logic runs once per goroutine regardless of how the
// pool schedules the two connections.
func TestPay_ConcurrentPaymentsOnlyOneSucceeds(t *testing.T) {
	db := testDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB(): %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	s := newOrderService(t, db)
	order := createPendingOrder(t, s, 1)
	all := &actions.DataPermission{DataScope: actions.DataScopeAll}

	var wg sync.WaitGroup
	errs := make([]error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = s.Pay(order.Id, all)
		}(i)
	}
	wg.Wait()

	successes, failures := 0, 0
	for _, err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrOrderNotPending):
			failures++
		default:
			t.Fatalf("unexpected error from a concurrent Pay: %v", err)
		}
	}
	if successes != 1 || failures != 1 {
		t.Fatalf("got %d successes and %d failures, want exactly 1 and 1", successes, failures)
	}
}

// -- data permission ---------------------------------------------------------

func TestGetPage_SelfScopeOnlySeesOwnOrders(t *testing.T) {
	enableDataPermission(t)
	db := testDB(t)
	s := newOrderService(t, db)

	createPendingOrder(t, s, 1) // belongs to user 1
	createPendingOrder(t, s, 2) // belongs to user 2

	var list []models.Order
	count, err := s.GetPage(&orderdto.OrderSearchReq{}, &actions.DataPermission{
		DataScope: actions.DataScopeSelf,
		UserId:    1,
	}, &list)
	if err != nil {
		t.Fatalf("GetPage: %v", err)
	}
	if count != 1 || len(list) != 1 {
		t.Fatalf("got %d orders, want exactly the 1 belonging to user 1", count)
	}
	if list[0].UserId != 1 {
		t.Errorf("returned order belongs to user %d, not the caller", list[0].UserId)
	}
}

func TestGetPage_AllScopeSeesEveryOrder(t *testing.T) {
	enableDataPermission(t)
	db := testDB(t)
	s := newOrderService(t, db)

	createPendingOrder(t, s, 1)
	createPendingOrder(t, s, 2)

	var list []models.Order
	count, err := s.GetPage(&orderdto.OrderSearchReq{}, &actions.DataPermission{DataScope: actions.DataScopeAll}, &list)
	if err != nil {
		t.Fatalf("GetPage: %v", err)
	}
	if count != 2 {
		t.Fatalf("got %d orders, want 2", count)
	}
}

// An invalid/unrecognized data_scope must fail closed - match nothing -
// never fall back to "see everything". This is core's own documented
// contract (contract/actions.Permission's default case), exercised here
// against app-order's own table to confirm the fail-closed behaviour
// actually reaches a hand-written Service's query, not just core's own
// unit tests.
func TestGetPage_InvalidScopeSeesNothing(t *testing.T) {
	enableDataPermission(t)
	db := testDB(t)
	s := newOrderService(t, db)

	createPendingOrder(t, s, 1)
	createPendingOrder(t, s, 2)

	var list []models.Order
	count, err := s.GetPage(&orderdto.OrderSearchReq{}, &actions.DataPermission{DataScope: "not-a-real-scope"}, &list)
	if err != nil {
		t.Fatalf("GetPage: %v", err)
	}
	if count != 0 || len(list) != 0 {
		t.Fatalf("an invalid data scope returned %d orders, want 0 (fail closed)", count)
	}
}

func TestGet_ReturnsOrderWithItemsPreloaded(t *testing.T) {
	enableDataPermission(t)
	db := testDB(t)
	s := newOrderService(t, db)

	created := createPendingOrder(t, s, 1)

	var got models.Order
	err := s.Get(created.Id, &actions.DataPermission{DataScope: actions.DataScopeSelf, UserId: 1}, &got)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("got %d items, want the 1 created with the order", len(got.Items))
	}
	if got.Items[0].ProductName != "widget" {
		t.Errorf("item ProductName = %q, want widget", got.Items[0].ProductName)
	}
}

func TestGet_ScopedOutOrderReportsNotFoundNotForbidden(t *testing.T) {
	enableDataPermission(t)
	db := testDB(t)
	s := newOrderService(t, db)

	other := createPendingOrder(t, s, 2)

	var got models.Order
	err := s.Get(other.Id, &actions.DataPermission{DataScope: actions.DataScopeSelf, UserId: 1}, &got)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Get on another user's order returned %v, want gorm.ErrRecordNotFound", err)
	}
}
