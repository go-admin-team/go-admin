package dto

import (
	"testing"

	"go-admin/app/demo/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

// 通用 Action 依赖 DTO 与 Model 实现一组接口。这些约束在编译期无法完全覆盖
// （接口是在路由注册处才被要求的），因此用测试锁定，避免改动后在运行时才暴露。

func TestImplementsIndexInterface(t *testing.T) {
	var _ dto.Index = (*DemoProductSearch)(nil)
}

func TestImplementsControlInterface(t *testing.T) {
	var _ dto.Control = (*DemoProductControl)(nil)
	var _ dto.Control = (*DemoProductById)(nil)
}

func TestModelImplementsActiveRecord(t *testing.T) {
	var _ common.ActiveRecord = (*models.DemoProduct)(nil)
}

// Generate 必须返回副本：通用 Action 在并发请求间复用同一个实例，
// 就地返回会导致请求之间串数据。
func TestGenerateReturnsCopy(t *testing.T) {
	src := &DemoProductControl{Id: 1, Name: "原始"}
	got := src.Generate().(*DemoProductControl)

	if got == src {
		t.Fatal("Generate 返回了同一指针，应返回副本")
	}
	got.Name = "被修改"
	if src.Name != "原始" {
		t.Errorf("修改副本影响了原对象：src.Name = %q", src.Name)
	}
}

func TestSearchGenerateReturnsCopy(t *testing.T) {
	src := &DemoProductSearch{Name: "原始"}
	got := src.Generate().(*DemoProductSearch)

	if got == src {
		t.Fatal("Generate 返回了同一指针，应返回副本")
	}
	got.Name = "被修改"
	if src.Name != "原始" {
		t.Errorf("修改副本影响了原对象：src.Name = %q", src.Name)
	}
}

func TestModelGenerateReturnsCopy(t *testing.T) {
	src := &models.DemoProduct{Name: "原始"}
	got := src.Generate().(*models.DemoProduct)

	if got == src {
		t.Fatal("Generate 返回了同一指针，应返回副本")
	}
	got.Name = "被修改"
	if src.Name != "原始" {
		t.Errorf("修改副本影响了原对象：src.Name = %q", src.Name)
	}
}

// GenerateM 组装落库对象，主键需正确传递，否则更新会退化成插入。
func TestGenerateMCarriesId(t *testing.T) {
	c := &DemoProductControl{Id: 42, Name: "示例", Code: "P-42", Price: 9.9}
	m, err := c.GenerateM()
	if err != nil {
		t.Fatalf("GenerateM 返回错误: %v", err)
	}
	p, ok := m.(*models.DemoProduct)
	if !ok {
		t.Fatalf("GenerateM 返回类型错误: %T", m)
	}
	if p.Id != 42 {
		t.Errorf("主键未传递: got %d, want 42", p.Id)
	}
	if p.Name != "示例" || p.Code != "P-42" || p.Price != 9.9 {
		t.Errorf("字段映射有误: %+v", p)
	}
}

func TestTableName(t *testing.T) {
	if got := (models.DemoProduct{}).TableName(); got != "demo_product" {
		t.Errorf("TableName() = %q, want %q", got, "demo_product")
	}
}
