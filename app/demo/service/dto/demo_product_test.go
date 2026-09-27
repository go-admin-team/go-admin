package dto

import (
	"testing"

	"go-admin/app/demo/models"
)

// ToModel builds the row to write; the key has to reach it, or an update
// matches no row.
func TestToModelCarriesId(t *testing.T) {
	c := &DemoProductControl{Id: 42, Name: "示例", Code: "P-42", Price: 9.9}
	p, err := c.ToModel()
	if err != nil {
		t.Fatalf("ToModel: %v", err)
	}
	if p.Id != 42 {
		t.Errorf("key not carried: got %d, want 42", p.Id)
	}
	if p.Name != "示例" || p.Code != "P-42" || p.Price != 9.9 {
		t.Errorf("fields mapped wrongly: %+v", p)
	}
}

func TestTableName(t *testing.T) {
	if got := (models.DemoProduct{}).TableName(); got != "demo_product" {
		t.Errorf("TableName() = %q, want %q", got, "demo_product")
	}
}
