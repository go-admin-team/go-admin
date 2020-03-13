package nacos_client

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"log"
	"os"
	"strconv"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-09 16:39
**/

type NacosClient struct {
	clientConfigValid  bool
	serverConfigsValid bool
	agent              http_agent.IHttpAgent
	clientConfig       constant.ClientConfig
	serverConfigs      []constant.ServerConfig
}

// 设置 clientConfig
func (client *NacosClient) SetClientConfig(config constant.ClientConfig) (err error) {
	if config.TimeoutMs <= 0 {
		err = errors.New("[client.SetClientConfig] config.TimeoutMs should > 0")
		return
	}
	if config.TimeoutMs >= config.ListenInterval {
		err = errors.New("[client.SetClientConfig] config.TimeoutMs should < config.ListenInterval")
		return
	}

	if config.BeatInterval <= 0 {
		config.BeatInterval = 5 * 1000
	}
	if config.ListenInterval < 10*1000 {
		config.ListenInterval = 10 * 1000
	}

	if config.UpdateThreadNum <= 0 {
		config.UpdateThreadNum = 20
	}

	if config.CacheDir == "" {
		config.CacheDir = utils.GetCurrentPath() + string(os.PathSeparator) + "cache"
	}
	if config.LogDir == "" {
		config.LogDir = utils.GetCurrentPath() + string(os.PathSeparator) + "log"
	}
	log.Printf("[INFO] logDir:<%s>   cacheDir:<%s>", config.LogDir, config.CacheDir)
	client.clientConfig = config
	client.clientConfigValid = true

	return
}

// 设置 serverConfigs
func (client *NacosClient) SetServerConfig(configs []constant.ServerConfig) (err error) {
	if len(configs) <= 0 {
		client.serverConfigsValid = true
		//err = errors.New("[client.SetServerConfig] configs can not be empty")
		return
	}

	for i := 0; i < len(configs); i++ {
		if len(configs[i].IpAddr) <= 0 || configs[i].Port <= 0 || configs[i].Port > 65535 {
			err = errors.New("[client.SetServerConfig] configs[" + strconv.Itoa(i) + "] is invalid")
			return
		}
		if len(configs[i].ContextPath) <= 0 {
			configs[i].ContextPath = constant.DEFAULT_CONTEXT_PATH
		}
	}
	client.serverConfigs = configs
	client.serverConfigsValid = true
	return
}

// 获取 clientConfig
func (client *NacosClient) GetClientConfig() (config constant.ClientConfig, err error) {
	config = client.clientConfig
	if !client.clientConfigValid {
		err = errors.New("[client.GetClientConfig] invalid client config")
	}
	return
}

// 获取serverConfigs
func (client *NacosClient) GetServerConfig() (configs []constant.ServerConfig, err error) {
	configs = client.serverConfigs
	if !client.serverConfigsValid {
		err = errors.New("[client.GetServerConfig] invalid server configs")
	}
	return
}

func (client *NacosClient) SetHttpAgent(agent http_agent.IHttpAgent) (err error) {
	if agent == nil {
		err = errors.New("[client.SetHttpAgent] http agent can not be nil")
	} else {
		client.agent = agent
	}
	return
}

func (client *NacosClient) GetHttpAgent() (agent http_agent.IHttpAgent, err error) {
	if client.agent == nil {
		err = errors.New("[client.GetHttpAgent] invalid http agent")
	} else {
		agent = client.agent
	}
	return
}
