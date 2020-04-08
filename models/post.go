package models

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type Post struct {
	//岗位编号
	PostId int64 `gorm:"column:post_id;primary_key" json:"postId" example:"1" extensions:"x-description=标示"`

	//岗位名称
	PostName string `gorm:"column:post_name" json:"postName"`

	//岗位代码
	PostCode string `gorm:"column:post_code" json:"postCode"`

	//岗位排序
	Sort int `gorm:"column:sort" json:"sort"`

	//状态
	Status string `gorm:"column:status" json:"status"`

	//描述
	Remark string `gorm:"column:remark" json:"remark"`

	//创建时间
	CreateTime string `gorm:"column:create_time" json:"createTime"`

	//最后修改时间
	UpdateTime string `gorm:"column:update_time" json:"updateTime"`

	//是否删除
	IsDel int `gorm:"column:is_del" json:"isDel"`

	CreateBy string `gorm:"column:create_by" json:"createBy"`

	UpdateBy string `gorm:"column:update_by" json:"updateBy"`

	DataScope string `gorm:"-" json:"dataScope"`

	Params string `gorm:"-" json:"params"`
}

func (e *Post) Create() (Post, error) {
	var doc Post
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	result := orm.Eloquent.Table("sys_post").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *Post) Get() (Post, error) {
	var doc Post

	table := orm.Eloquent.Table("sys_post")
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

	if err := table.Where("is_del = 0").First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Post) GetList() ([]Post, error) {
	var doc []Post

	table := orm.Eloquent.Table("sys_post")
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

	if err := table.Where("is_del = 0").Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Post) GetPage(pageSize int, pageIndex int) ([]Post, int32, error) {
	var doc []Post

	table := orm.Eloquent.Select("*").Table("sys_post")
	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}
	if e.PostName != "" {
		table = table.Where("post_name = ?", e.PostName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = utils.StringToInt64(e.DataScope)
	table = dataPermission.GetDataScope("sys_post", table)

	var count int32

	if err := table.Where("is_del = 0").Order("sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("is_del = 0").Count(&count)
	return doc, count, nil
}

func (e *Post) Update(id int64) (update Post, err error) {
	e.UpdateTime = utils.GetCurrntTime()
	if err = orm.Eloquent.Table("sys_post").First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_post").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Post) Delete(id int64) (success bool, err error) {
	var mp = map[string]string{}
	mp["is_del"] = "1"
	mp["update_time"] = utils.GetCurrntTime()
	mp["update_by"] = e.UpdateBy
	if err = orm.Eloquent.Table("sys_post").Where("post_id = ?", id).Update(mp).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
