package vo

import "github.com/nacos-group/nacos-sdk-go/model"

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-09 10:03
**/

type RegisterInstanceParam struct {
	Ip          string            `param:"ip"`
	Port        uint64            `param:"port"`
	Tenant      string            `param:"tenant"`
	Weight      float64           `param:"weight"`
	Enable      bool              `param:"enabled"`
	Healthy     bool              `param:"healthy"`
	Metadata    map[string]string `param:"metadata"`
	ClusterName string            `param:"clusterName"`
	ServiceName string            `param:"serviceName"`
	GroupName   string            `param:"groupName"`
	Ephemeral   bool              `param:"ephemeral"`
}

type DeregisterInstanceParam struct {
	Ip          string `param:"ip"`
	Port        uint64 `param:"port"`
	Tenant      string `param:"tenant"`
	Cluster     string `param:"cluster"`
	ServiceName string `param:"serviceName"`
	GroupName   string `param:"groupName"`
	Ephemeral   bool   `param:"ephemeral"`
}

type GetServiceParam struct {
	Clusters    []string `param:"clusters"`
	ServiceName string   `param:"serviceName"`
	GroupName   string   `param:"groupName"`
}

type GetAllServiceInfoParam struct {
	Clusters  []string `param:"clusters"`
	NameSpace string   `param:"nameSpace"`
	GroupName string   `param:"groupName"`
}

type GetServiceListParam struct {
	StartPage   uint32 `param:"startPg"`
	PageSize    uint32 `param:"pgSize"`
	Keyword     string `param:"keyword"`
	NamespaceId string `param:"namespaceId"`
}

type GetServiceInstanceParam struct {
	Tenant      string `param:"tenant"`
	HealthyOnly bool   `param:"healthyOnly"`
	Cluster     string `param:"cluster"`
	ServiceName string `param:"serviceName"`
	Ip          string `param:"ip"`
	Port        uint64 `param:"port"`
}

type GetServiceDetailParam struct {
	ServiceName string `param:"serviceName"`
}

type SubscribeParam struct {
	ServiceName       string   `param:"serviceName"`
	Clusters          []string `param:"clusters"`
	GroupName         string   `param:"groupName"`
	SubscribeCallback func(services []model.SubscribeService, err error)
}

type SelectAllInstancesParam struct {
	Clusters    []string `param:"clusters"`
	ServiceName string   `param:"serviceName"`
	GroupName   string   `param:"groupName"`
}

type SelectInstancesParam struct {
	Clusters    []string `param:"clusters"`
	ServiceName string   `param:"serviceName"`
	GroupName   string   `param:"groupName"`
	HealthyOnly bool     `param:"healthyOnly"`
}

type SelectOneHealthInstanceParam struct {
	Clusters    []string `param:"clusters"`
	ServiceName string   `param:"serviceName"`
	GroupName   string   `param:"groupName"`
}
