package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/nacos-group/nacos-sdk-go/common/util"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func GetFileName(cacheKey string, cacheDir string) string {
	return cacheDir + string(os.PathSeparator) + cacheKey
}

func WriteServicesToFile(service model.Service, cacheDir string) {
	util.MkdirIfNecessary(cacheDir)
	sb, _ := json.Marshal(service)
	domFileName := GetFileName(utils.GetServiceCacheKey(service.Name, service.Clusters), cacheDir)

	err := ioutil.WriteFile(domFileName, sb, 0666)
	if err != nil {
		log.Printf("[ERROR]:faild to write name cache:%s ,value:%s ,err:%s \n", domFileName, string(sb), err.Error())
	}

}

func ReadServicesFromFile(cacheDir string) map[string]model.Service {
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		log.Printf("[ERROR]:read cacheDir:%s failed!err:%s \n", cacheDir, err.Error())
		return nil
	}
	serviceMap := map[string]model.Service{}
	for _, f := range files {
		fileName := GetFileName(f.Name(), cacheDir)
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Printf("[ERROR]:failed to read name cache file:%s,err:%s! ", fileName, err.Error())
			continue
		}

		s := string(b)
		service := utils.JsonToService(s)

		if service == nil {
			continue
		}

		serviceMap[f.Name()] = *service
	}

	log.Printf("finish loading name cache, total: " + strconv.Itoa(len(files)))
	return serviceMap
}

func WriteConfigToFile(cacheKey string, cacheDir string, content string) {
	util.MkdirIfNecessary(cacheDir)
	fileName := GetFileName(cacheKey, cacheDir)
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		log.Printf("[ERROR]:faild to write config  cache:%s ,value:%s ,err:%s \n", fileName, string(content), err.Error())
	}
}

func ReadConfigFromFile(cacheKey string, cacheDir string) (string, error) {
	fileName := GetFileName(cacheKey, cacheDir)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to read config cache file:%s,err:%s! ", fileName, err.Error()))
	}
	return string(b), nil
}
