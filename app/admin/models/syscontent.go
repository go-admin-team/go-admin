package models

import (
	orm "go-admin/common/global"
	"go-admin/tools"
)

type SysContent struct {
	Id        int    `json:"id" gorm:"type:int(11);primary_key;auto_increment"` // id
	CateId    string `json:"cateId" gorm:"type:int(11);"`                       // 分类id
	Name      string `json:"name" gorm:"type:varchar(255);"`                    // 名称
	Status    string `json:"status" gorm:"type:int(1);"`                        // 状态
	Img       string `json:"img" gorm:"type:varchar(255);"`                     // 图片
	Content   string `json:"content" gorm:"type:text;"`                         // 内容
	Remark    string `json:"remark" gorm:"type:varchar(255);"`                  // 备注
	Sort      string `json:"sort" gorm:"type:int(4);"`                          // 排序
	CreateBy  string `json:"createBy" gorm:"type:varchar(128);"`                // 创建者
	UpdateBy  string `json:"updateBy" gorm:"type:varchar(128);"`                // 更新者
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
	BaseModel
}

func (SysContent) TableName() string {
	return "sys_content"
}

// 创建SysContent
func (e *SysContent) Create() (SysContent, error) {
	var doc SysContent
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取SysContent
func (e *SysContent) Get() (SysContent, error) {
	var doc SysContent
	table := orm.Eloquent.Table(e.TableName())

	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	if e.CateId != "" {
		table = table.Where("cate_id = ?", e.CateId)
	}

	if e.Name != "" {
		table = table.Where("name like ?", "%"+e.Name+"%")
	}

	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取SysContent带分页
func (e *SysContent) GetPage(pageSize int, pageIndex int) ([]SysContent, int, error) {
	var doc []SysContent

	table := orm.Eloquent.Table(e.TableName())

	if e.CateId != "" {
		table = table.Where("cate_id = ?", e.CateId)
	}

	if e.Name != "" {
		table = table.Where("name like ?", "%"+e.Name+"%")
	}

	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}
	var count int64

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	//table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, int(count), nil
}

// 更新SysContent
func (e *SysContent) Update(id int) (update SysContent, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除SysContent
func (e *SysContent) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&SysContent{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysContent) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&SysContent{}).Error; err != nil {
		return
	}
	Result = true
	return
}
