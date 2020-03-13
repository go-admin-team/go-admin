package config_client

import (
	"github.com/nacos-group/nacos-sdk-go/vo"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-08 10:03
**/

//go:generate mockgen -destination ../../mock/mock_config_client_interface.go -package mock -source=./config_client_interface.go

type IConfigClient interface {
	// 获取配置
	// dataId  require
	// group   require
	// tenant ==>nacos.namespace optional
	GetConfig(param vo.ConfigParam) (string, error)

	// 发布配置
	// dataId  require
	// group   require
	// content require
	// tenant ==>nacos.namespace optional
	PublishConfig(param vo.ConfigParam) (bool, error)

	// 删除配置
	// dataId  require
	// group   require
	// tenant ==>nacos.namespace optional
	DeleteConfig(param vo.ConfigParam) (bool, error)

	// 监听配置
	// dataId  require
	// group   require
	// tenant ==>nacos.namespace optional
	ListenConfig(params vo.ConfigParam) (err error)
}
