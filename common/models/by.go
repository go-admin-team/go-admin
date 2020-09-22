package models

type ControlBy struct {
	CreateBy uint `gorm:"index;comment:'创建者'"`
	UpdateBy uint `gorm:"index;comment:'更新者'"`
}

func (e *ControlBy) SetCreateBy(createBy uint) {
	e.CreateBy = createBy
}

func (e *ControlBy) SetUpdateBy(updateBy uint) {
	e.UpdateBy = updateBy
}
