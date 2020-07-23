package models

import (
	orm "go-admin/global"
	"go-admin/tools"
)

type Post struct {
	PostId    int    `gorm:"primary_key;AUTO_INCREMENT" json:"postId"` //岗位编号
	PostName  string `gorm:"size:128;" json:"postName"`                //岗位名称
	PostCode  string `gorm:"size:128;" json:"postCode"`                //岗位代码
	Sort      int    `gorm:"" json:"sort"`                             //岗位排序
	Status    string `gorm:"size:4;" json:"status"`                    //状态
	Remark    string `gorm:"size:255;" json:"remark"`                  //描述
	CreateBy  string `gorm:"size:128;" json:"createBy"`
	UpdateBy  string `gorm:"size:128;" json:"updateBy"`
	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
	BaseModel
}

func (Post) TableName() string {
	return "sys_post"
}

func (e *Post) Create() (Post, error) {
	var doc Post
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *Post) Get() (Post, error) {
	var doc Post

	table := orm.Eloquent.Table(e.TableName())
	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}
	if e.PostName != "" {
		table = table.Where("post_name = ?", e.PostName)
	}
	if e.PostCode != "" {
		table = table.Where("post_code = ?", e.PostCode)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Post) GetList() ([]Post, error) {
	var doc []Post

	table := orm.Eloquent.Table(e.TableName())
	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}
	if e.PostName != "" {
		table = table.Where("post_name = ?", e.PostName)
	}
	if e.PostCode != "" {
		table = table.Where("post_code = ?", e.PostCode)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Post) GetPage(pageSize int, pageIndex int) ([]Post, int, error) {
	var doc []Post

	table := orm.Eloquent.Select("*").Table(e.TableName())
	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}
	if e.PostName != "" {
		table = table.Where("post_name = ?", e.PostName)
	}
	if e.PostCode != "" {
		table = table.Where("post_code = ?", e.PostCode)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("sys_post", table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Order("sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

func (e *Post) Update(id int) (update Post, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Post) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("post_id = ?", id).Delete(&Post{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *Post) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("post_id in (?)", id).Delete(&Post{}).Error; err != nil {
		return
	}
	Result = true
	return
}
