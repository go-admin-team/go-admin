package nacos_client

import (
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-09 16:32
**/

//go:generate mockgen -destination mock_nacos_client_interface.go -package nacos_client -source=./nacos_client_interface.go

type INacosClient interface {
	SetClientConfig(constant.ClientConfig) error
	SetServerConfig([]constant.ServerConfig) error
	GetClientConfig() (constant.ClientConfig, error)
	GetServerConfig() ([]constant.ServerConfig, error)
	SetHttpAgent(http_agent.IHttpAgent) error
	GetHttpAgent() (http_agent.IHttpAgent, error)
}
