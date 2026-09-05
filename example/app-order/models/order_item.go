package models

// orderItemTable mirrors orderTable's naming rationale: not a reserved word,
// and namespaced under app_ so a host scanning its schema can tell at a
// glance which tables an installed app owns.
const orderItemTable = "app_order_item"

// OrderItem is one line item of an Order. It carries no ControlBy of its
// own: data-scope is enforced once, on the parent Order, and an item is
// never queried on its own outside that parent (see service.Order.Get's
// Preload).
//
// OrderId+ProductName is unique on purpose, not just to have some index: it
// is what OrderService_test.go's mid-transaction-failure test relies on to
// force a real constraint violation after the parent Order row has already
// been inserted in the same transaction, proving the rollback actually
// undoes both writes rather than leaving the Order behind.
type OrderItem struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement;comment:primary key"`
	OrderId     int    `json:"orderId" gorm:"uniqueIndex:uk_app_order_item_product;comment:parent order id"`
	ProductName string `json:"productName" gorm:"type:varchar(255);uniqueIndex:uk_app_order_item_product;comment:product name"`
	Quantity    int    `json:"quantity" gorm:"comment:quantity"`
	PriceCents  int64  `json:"priceCents" gorm:"comment:unit price in cents"`
}

// TableName pins the row model to app_order_item; see orderItemTable.
func (OrderItem) TableName() string {
	return orderItemTable
}
