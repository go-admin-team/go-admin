package config_client

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/common/nacos_error"
	"github.com/nacos-group/nacos-sdk-go/common/util"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ConfigClient struct {
	nacos_client.INacosClient
	kmsClient      *kms.Client
	localConfigs   []vo.ConfigParam
	mutex          sync.Mutex
	configProxy    ConfigProxy
	configCacheDir string
}

func NewConfigClient(nc nacos_client.INacosClient) (ConfigClient, error) {
	config := ConfigClient{}
	config.INacosClient = nc
	clientConfig, err := nc.GetClientConfig()
	if err != nil {
		return config, err
	}
	serverConfig, err := nc.GetServerConfig()
	if err != nil {
		return config, err
	}
	httpAgent, err := nc.GetHttpAgent()
	if err != nil {
		return config, err
	}
	err = logger.InitLog(clientConfig.LogDir)
	if err != nil {
		return config, err
	}
	config.configCacheDir = clientConfig.CacheDir + string(os.PathSeparator) + "config"
	config.configProxy, err = NewConfigProxy(serverConfig, clientConfig, httpAgent)
	if clientConfig.OpenKMS {
		kmsClient, err := kms.NewClientWithAccessKey(clientConfig.RegionId, clientConfig.AccessKey, clientConfig.SecretKey)
		if err != nil {
			return config, err
		}
		config.kmsClient = kmsClient
	}

	return config, err
}

func (client *ConfigClient) sync() (clientConfig constant.ClientConfig,
	serverConfigs []constant.ServerConfig, agent http_agent.IHttpAgent, err error) {
	clientConfig, err = client.GetClientConfig()
	if err != nil {
		log.Println(err, ";do you call client.SetClientConfig()?")
	}
	if err == nil {
		serverConfigs, err = client.GetServerConfig()
		if err != nil {
			log.Println(err, ";do you call client.SetServerConfig()?")
		}
	}
	if err == nil {
		agent, err = client.GetHttpAgent()
		if err != nil {
			log.Println(err, ";do you call client.SetHttpAgent()?")
		}
	}
	return
}

func (client *ConfigClient) GetConfig(param vo.ConfigParam) (content string, err error) {
	content, err = client.getConfigInner(param)

	if err != nil {
		return "", err
	}

	return client.decrypt(param.DataId, content)
}

func (client *ConfigClient) decrypt(dataId, content string) (string, error) {
	if strings.HasPrefix(dataId, "cipher-") && client.kmsClient != nil {
		request := kms.CreateDecryptRequest()
		request.Method = "POST"
		request.Scheme = "https"
		request.AcceptFormat = "json"
		request.CiphertextBlob = content
		response, err := client.kmsClient.Decrypt(request)
		if err != nil {
			return "", errors.New("ksm decrypt failed")
		}
		content = response.Plaintext
	}

	return content, nil
}

func (client *ConfigClient) getConfigInner(param vo.ConfigParam) (content string, err error) {
	if len(param.DataId) <= 0 {
		err = errors.New("[client.GetConfig] param.dataId can not be empty")
	}
	if len(param.Group) <= 0 {
		err = errors.New("[client.GetConfig] param.group can not be empty")
	}
	clientConfig, _ := client.GetClientConfig()
	cacheKey := utils.GetConfigCacheKey(param.DataId, param.Group, clientConfig.NamespaceId)
	content, err = client.configProxy.GetConfigProxy(param, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)

	if err != nil {
		log.Printf("[ERROR] get config from server error:%s ", err.Error())
		if _, ok := err.(*nacos_error.NacosError); ok {
			nacosErr := err.(*nacos_error.NacosError)
			if nacosErr.ErrorCode() == "404" {
				cache.WriteConfigToFile(cacheKey, client.configCacheDir, "")
				return "", errors.New("config not found")
			}
			if nacosErr.ErrorCode() == "403" {
				return "", errors.New("get config forbidden")
			}
		}
		content, err = cache.ReadConfigFromFile(cacheKey, client.configCacheDir)
		if err != nil {
			log.Printf("[ERROR] get config from cache  error:%s ", err.Error())
			return "", errors.New("read config from both server and cache fail")
		}

	} else {
		cache.WriteConfigToFile(cacheKey, client.configCacheDir, content)
	}
	return content, nil
}

func (client *ConfigClient) PublishConfig(param vo.ConfigParam) (published bool,
	err error) {
	if len(param.DataId) <= 0 {
		err = errors.New("[client.PublishConfig] param.dataId can not be empty")
	}
	if len(param.Group) <= 0 {
		err = errors.New("[client.PublishConfig] param.group can not be empty")
	}
	if len(param.Content) <= 0 {
		err = errors.New("[client.PublishConfig] param.content can not be empty")
	}
	clientConfig, _ := client.GetClientConfig()
	return client.configProxy.PublishConfigProxy(param, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)
}

func (client *ConfigClient) DeleteConfig(param vo.ConfigParam) (deleted bool,
	err error) {
	if len(param.DataId) <= 0 {
		err = errors.New("[client.DeleteConfig] param.dataId can not be empty")
	}
	if len(param.Group) <= 0 {
		err = errors.New("[client.DeleteConfig] param.group can not be empty")
	}

	clientConfig, _ := client.GetClientConfig()
	return client.configProxy.DeleteConfigProxy(param, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)
}

func (client *ConfigClient) AddConfigToListen(params []vo.ConfigParam) (err error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	var newParams []vo.ConfigParam
	// 去掉重复添加的
	for _, newParam := range params {
		flag := true
		for _, param := range client.localConfigs {
			if param.Group == newParam.Group && param.DataId == newParam.DataId {
				flag = false
				break
			}
		}
		if flag {
			newParams = append(newParams, newParam)
		}
	}
	client.localConfigs = append(client.localConfigs, newParams...)
	return
}

func (client *ConfigClient) ListenConfig(param vo.ConfigParam) (err error) {
	go func() {
		for {
			clientConfig, serverConfigs, agent, err := client.sync()
			// 创建计时器
			var timer *time.Timer
			if err == nil {
				timer = time.NewTimer(time.Duration(clientConfig.ListenInterval) * time.Millisecond)
			}
			client.listenConfigTask(clientConfig, serverConfigs, agent, param)
			<-timer.C
		}
	}()

	return nil
}

func (client *ConfigClient) listenConfigTask(clientConfig constant.ClientConfig,
	serverConfigs []constant.ServerConfig, agent http_agent.IHttpAgent, param vo.ConfigParam) {
	var listeningConfigs string
	// 检查&拼接监听参数
	client.mutex.Lock()
	if len(param.DataId) <= 0 {
		log.Fatalf("[client.ListenConfig] DataId can not be empty")
		return
	}
	if len(param.Group) <= 0 {
		log.Fatalf("[client.ListenConfig] Group can not be empty")
		return
	}
	var tenant string
	if len(clientConfig.NamespaceId) > 0 {
		tenant = clientConfig.NamespaceId
	}

	for _, config := range client.localConfigs {
		if config.Group == param.Group && config.DataId == param.DataId {
			param.Content = config.Content
			break
		}
	}

	var md5 string
	if len(param.Content) > 0 {
		md5 = util.Md5(param.Content)
	}
	if len(tenant) > 0 {
		listeningConfigs += param.DataId + constant.SPLIT_CONFIG_INNER + param.Group + constant.SPLIT_CONFIG_INNER +
			md5 + constant.SPLIT_CONFIG_INNER + tenant + constant.SPLIT_CONFIG
	} else {
		listeningConfigs += param.DataId + constant.SPLIT_CONFIG_INNER + param.Group + constant.SPLIT_CONFIG_INNER +
			md5 + constant.SPLIT_CONFIG
	}

	client.mutex.Unlock()
	// http 请求
	params := make(map[string]string)
	params[constant.KEY_LISTEN_CONFIGS] = listeningConfigs
	var changed string

	for _, serverConfig := range client.configProxy.GetServerList() {
		path := client.buildBasePath(serverConfig) + "/listener"
		changedTmp, err := listen(agent, path, clientConfig.TimeoutMs, clientConfig.ListenInterval, params)
		if err == nil {
			changed = changedTmp
			break
		} else {
			if _, ok := err.(*nacos_error.NacosError); ok {
				changed = changedTmp
				break
			} else {
				log.Println("[client.ListenConfig] listen config error:", err.Error())
			}
		}
	}

	if strings.ToLower(strings.Trim(changed, " ")) == "" {
		log.Println("[client.ListenConfig] no change")
	} else {
		log.Print("[client.ListenConfig] config changed:" + changed)
		client.updateLocalConfig(changed, param)
	}
}

func listen(agent http_agent.IHttpAgent, path string,
	timeoutMs uint64, listenInterval uint64,
	params map[string]string) (changed string, err error) {
	header := map[string][]string{
		"Content-Type":         {"application/x-www-form-urlencoded"},
		"Long-Pulling-Timeout": {strconv.FormatUint(listenInterval, 10)},
	}
	log.Println("[client.ListenConfig] request url:", path, " ;params:", params, " ;header:", header)
	var response *http.Response
	response, err = agent.Post(path, header, timeoutMs, params)
	if err == nil {
		bytes, errRead := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if errRead != nil {
			err = errRead
		} else {
			if response.StatusCode == 200 {
				changed = string(bytes)
			} else {
				err = nacos_error.NewNacosError(strconv.Itoa(response.StatusCode), string(bytes), nil)
			}
		}
	}
	return
}

func (client *ConfigClient) updateLocalConfig(changed string, param vo.ConfigParam) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	changedConfigs := strings.Split(changed, "%01")
	for _, config := range changedConfigs {
		attrs := strings.Split(config, "%02")
		if len(attrs) == 2 {
			content, err := client.getConfigInner(vo.ConfigParam{
				DataId: attrs[0],
				Group:  attrs[1],
			})
			if err != nil {
				log.Println("[client.updateLocalConfig] update config failed:", err.Error())
			} else {
				client.putLocalConfig(vo.ConfigParam{
					DataId:  attrs[0],
					Group:   attrs[1],
					Content: content,
				})

				// call listener:
				decrept, _ := client.decrypt(attrs[0], content)
				param.OnChange("", attrs[1], attrs[0], decrept)

			}
		} else if len(attrs) == 3 {
			content, err := client.getConfigInner(vo.ConfigParam{
				DataId: attrs[0],
				Group:  attrs[1],
			})
			if err != nil {
				log.Println("[client.updateLocalConfig] update config failed:", err.Error())
			} else {
				client.putLocalConfig(vo.ConfigParam{
					DataId:  attrs[0],
					Group:   attrs[1],
					Content: content,
				})

				// call listener:
				decrept, _ := client.decrypt(attrs[0], content)
				param.OnChange(attrs[2], attrs[1], attrs[0], decrept)
			}
		}
	}
	log.Println("[client.updateLocalConfig] update config complete")
	log.Println("[client.localConfig] ", client.localConfigs)
}

func (client *ConfigClient) putLocalConfig(config vo.ConfigParam) {
	if len(config.DataId) > 0 && len(config.Group) > 0 {
		exist := false
		for i := 0; i < len(client.localConfigs); i++ {
			if len(client.localConfigs[i].DataId) > 0 && len(client.localConfigs[i].Group) > 0 &&
				config.DataId == client.localConfigs[i].DataId && config.Group == client.localConfigs[i].Group {
				// 本地存在 则更新
				client.localConfigs[i] = config
				exist = true
				break
			}
		}
		if !exist {
			// 本地不存在 放入
			client.localConfigs = append(client.localConfigs, config)
		}
	}
	log.Println("[client.putLocalConfig] putLocalConfig success")
}

func (client *ConfigClient) buildBasePath(serverConfig constant.ServerConfig) (basePath string) {
	basePath = "http://" + serverConfig.IpAddr + ":" +
		strconv.FormatUint(serverConfig.Port, 10) + serverConfig.ContextPath + constant.CONFIG_PATH
	return
}
