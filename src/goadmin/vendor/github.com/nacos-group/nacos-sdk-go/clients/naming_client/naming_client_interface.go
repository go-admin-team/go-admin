package naming_client

import (
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-09 09:56
**/

//go:generate mockgen -destination ../../mock/mock_service_client_interface.go -package mock -source=./service_client_interface.go

type INamingClient interface {
	// 注册服务实例
	RegisterInstance(param vo.RegisterInstanceParam) (bool, error)
	// 注销服务实例
	DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error)
	// 获取服务信息
	GetService(param vo.GetServiceParam) (model.Service, error)
	//获取所有的实例列表
	SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error)
	// 获取实例列表
	SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error)
	//获取一个健康的实例
	SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error)
	// 服务监听
	Subscribe(param *vo.SubscribeParam) error
	//取消监听
	Unsubscribe(param *vo.SubscribeParam) error

	//获取全部服务信息
	GetAllServicesInfo(param vo.GetAllServiceInfoParam) ([]model.Service, error)
}
