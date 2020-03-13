package naming_client

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"log"
)

type SubscribeCallback struct {
	callbackFuncsMap cache.ConcurrentMap
}

func NewSubscribeCallback() SubscribeCallback {
	ed := SubscribeCallback{}
	ed.callbackFuncsMap = cache.NewConcurrentMap()
	return ed
}

func (ed *SubscribeCallback) AddCallbackFuncs(serviceName string, clusters string, callbackFunc *func(services []model.SubscribeService, err error)) {
	log.Printf("[INFO] adding " + serviceName + " with " + clusters + " to listener map")
	key := utils.GetServiceCacheKey(serviceName, clusters)
	var funcs []*func(services []model.SubscribeService, err error)
	old, ok := ed.callbackFuncsMap.Get(key)
	if ok {
		funcs = append(funcs, old.([]*func(services []model.SubscribeService, err error))...)
	}
	funcs = append(funcs, callbackFunc)
	ed.callbackFuncsMap.Set(key, funcs)
}

func (ed *SubscribeCallback) RemoveCallbackFuncs(serviceName string, clusters string, callbackFunc *func(services []model.SubscribeService, err error)) {
	log.Printf("[INFO] removing " + serviceName + " with " + clusters + " to listener map")
	key := utils.GetServiceCacheKey(serviceName, clusters)
	funcs, ok := ed.callbackFuncsMap.Get(key)
	if ok && funcs != nil {
		var newFuncs []*func(services []model.SubscribeService, err error)
		for _, funcItem := range funcs.([]*func(services []model.SubscribeService, err error)) {
			if funcItem != callbackFunc {
				newFuncs = append(newFuncs, funcItem)
			}
		}
		ed.callbackFuncsMap.Set(key, newFuncs)
	}

}

func (ed *SubscribeCallback) ServiceChanged(service *model.Service) {
	if service == nil || service.Name == "" {
		return
	}
	key := utils.GetServiceCacheKey(service.Name, service.Clusters)
	funcs, ok := ed.callbackFuncsMap.Get(key)
	if ok {
		for _, funcItem := range funcs.([]*func(services []model.SubscribeService, err error)) {
			var subscribeServices []model.SubscribeService
			if len(service.Hosts) == 0 {
				(*funcItem)(subscribeServices, errors.New("[client.Subscribe] subscribe failed,hosts is empty"))
				return
			}
			for _, host := range service.Hosts {
				var subscribeService model.SubscribeService
				subscribeService.Valid = host.Valid
				subscribeService.Port = host.Port
				subscribeService.Ip = host.Ip
				subscribeService.Metadata = host.Metadata
				subscribeService.ServiceName = host.ServiceName
				subscribeService.ClusterName = host.ClusterName
				subscribeService.Weight = host.Weight
				subscribeService.InstanceId = host.InstanceId
				subscribeService.Enable = host.Enable
				subscribeServices = append(subscribeServices, subscribeService)
			}
			(*funcItem)(subscribeServices, nil)
		}
	}
}
