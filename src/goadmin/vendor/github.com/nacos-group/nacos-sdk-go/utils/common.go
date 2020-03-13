package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func CurrentMillis() int64 {
	return time.Now().UnixNano() / 1e6
}

var (
	GZIP_MAGIC = []byte("\x1F\x8B")
)

func TryDecompressData(data []byte) string {

	if !IsGzipFile(data) {
		//fmt.Println("data format: plain text")
		return string(data)
	}

	//fmt.Println("data format: gzip")

	reader, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		log.Printf("[ERROR]:failed to decompress gzip data,err:%s \n", err.Error())
		return ""
	}

	defer reader.Close()
	bs, err1 := ioutil.ReadAll(reader)

	if err1 != nil {
		log.Printf("[ERROR]:failed to decompress gzip data,err:%s \n", err1.Error())
		return ""
	}

	return string(bs)
}

func IsGzipFile(data []byte) bool {
	if len(data) < 2 {
		return false
	}

	return bytes.HasPrefix(data, GZIP_MAGIC)
}

func JsonToService(result string) *model.Service {
	var service model.Service
	err := json.Unmarshal([]byte(result), &service)
	if err != nil {
		log.Printf("[ERROR]:failed to unmarshal json string:%s err:%v \n", result, err.Error())
		return nil
	}
	if len(service.Hosts) == 0 {
		log.Printf("[WARN]:instance list is empty,json string:%s \n", result)
		return nil
	}
	return &service

}
func ToJsonString(object interface{}) string {
	js, _ := json.Marshal(object)
	return string(js)
}

func GetCurrentPath() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("can not get current path, exit!!")
		os.Exit(1)
	}

	return dir
}

func GetGroupName(serviceName string, groupName string) string {
	return groupName + constant.SERVICE_INFO_SPLITER + serviceName
}

func GetServiceCacheKey(serviceName string, clusters string) string {
	if clusters == "" {
		return serviceName
	}
	return serviceName + constant.SERVICE_INFO_SPLITER + clusters
}

func GetConfigCacheKey(dataId string, group string, tenant string) string {
	return dataId + constant.CONFIG_INFO_SPLITER + group + constant.CONFIG_INFO_SPLITER + tenant
}

var localIP = ""

func LocalIP() string {
	if localIP == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.Printf("[ERROR]:get InterfaceAddres failed,err:%s \n", err.Error())
			return ""
		}
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					localIP = ipnet.IP.String()
					log.Printf("InitLocalIp, LocalIp:%s \n", localIP)
					break
				}
			}
		}
	}
	return localIP
}

func GetDurationWithDefault(metadata map[string]string, key string, defaultDuration time.Duration) time.Duration {
	data, ok := metadata[key]
	if ok {
		value, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			log.Printf("key:%s is not a number \n", key)
			return defaultDuration
		}
		return time.Duration(value)
	}
	return defaultDuration
}
