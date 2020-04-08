package models

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type Config struct {

	//编码
	ConfigId int64 `json:"configId" gorm:"column:config_id;primary_key"`
	//参数名称
	ConfigName string `json:"configName" gorm:"column:config_name"`
	//参数键名
	ConfigKey string `json:"configKey" gorm:"column:config_key"`
	//参数键值
	ConfigValue string `json:"configValue" gorm:"column:config_value"`
	//是否系统内置
	ConfigType string `json:"configType" gorm:"column:config_type"`
	//备注
	Remark string `json:"remark" gorm:"column:remark"`

	CreateBy   string `json:"createBy" gorm:"column:create_by"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`
	UpdateTime string `json:"updateTime" gorm:"column:update_time"`
	DataScope  string `json:"dataScope" gorm:"-"`
	Params     string `json:"params" gorm:"column:params"`
	IsDel      int    `json:"isDel" gorm:"column:is_del"`
}

// Config 创建
func (e *Config) Create() (Config, error) {
	var doc Config
	doc.IsDel = 0
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	result := orm.Eloquent.Table("sys_config").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取 Config
func (e *Config) Get() (Config, error) {
	var doc Config

	table := orm.Eloquent.Table("sys_config")
	if e.ConfigId != 0 {
		table = table.Where("config_id = ?", e.ConfigId)
	}

	if e.ConfigKey != "" {
		table = table.Where("config_key = ?", e.ConfigKey)
	}

	if err := table.Where("is_del = 0").First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Config) GetPage(pageSize int, pageIndex int) ([]Config, int32, error) {
	var doc []Config

	table := orm.Eloquent.Select("*").Table("sys_config")

	if e.ConfigName != "" {
		table = table.Where("config_name = ?", e.ConfigName)
	}
	if e.ConfigKey != "" {
		table = table.Where("config_key = ?", e.ConfigKey)
	}
	if e.ConfigType != "" {
		table = table.Where("config_type = ?", e.ConfigType)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = utils.StringToInt64(e.DataScope)
	table = dataPermission.GetDataScope("sys_config", table)

	var count int32

	if err := table.Where("is_del = 0").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("is_del = 0").Count(&count)
	return doc, count, nil
}

func (e *Config) Update(id int64) (update Config, err error) {
	e.UpdateTime = utils.GetCurrntTime()
	if err = orm.Eloquent.Table("sys_config").Where("config_id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_config").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Config) Delete(id int64) (success bool, err error) {
	var mp = map[string]string{}
	mp["is_del"] = "1"
	mp["update_time"] = utils.GetCurrntTime()
	mp["update_by"] = e.UpdateBy
	if err = orm.Eloquent.Table("sys_config").Where("config_id = ?", id).Update(mp).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
