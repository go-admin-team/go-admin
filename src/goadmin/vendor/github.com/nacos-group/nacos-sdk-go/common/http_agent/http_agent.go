package http_agent

import (
	"github.com/go-errors/errors"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"io/ioutil"
	"log"
	"net/http"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-10 11:26
**/
type HttpAgent struct {
}

func (agent *HttpAgent) Get(path string, header http.Header, timeoutMs uint64,
	params map[string]string) (response *http.Response, err error) {
	return get(path, header, timeoutMs, params)
}

func (agent *HttpAgent) RequestOnlyResult(method string, path string, header http.Header, timeoutMs uint64, params map[string]string) string {
	var response *http.Response
	var err error
	switch method {
	case http.MethodGet:
		response, err = agent.Get(path, header, timeoutMs, params)
		break
	case http.MethodPost:
		response, err = agent.Post(path, header, timeoutMs, params)
		break
	case http.MethodPut:
		response, err = agent.Put(path, header, timeoutMs, params)
		break
	case http.MethodDelete:
		response, err = agent.Delete(path, header, timeoutMs, params)
		break
	default:
		log.Printf("[ERROR]:request method[%s], path[%s],header:[%s],params:[%s], not avaliable method ", method, path, utils.ToJsonString(header), utils.ToJsonString(params))
	}
	if err != nil {
		log.Printf("[ERROR]:request method[%s],request path[%s],header:[%s],params:[%s],err:%s", method, path, utils.ToJsonString(header), utils.ToJsonString(params), err.Error())
		return ""
	}
	if response.StatusCode != 200 {
		log.Printf("[ERROR]:request method[%s],request path[%s],header:[%s],params:[%s],status code error:%d", method, path, utils.ToJsonString(header), utils.ToJsonString(params), response.StatusCode)
		return ""
	}
	bytes, errRead := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if errRead != nil {
		log.Printf("[ERROR]:request method[%s],request path[%s],header:[%s],params:[%s],read error:%s", method, path, utils.ToJsonString(header), utils.ToJsonString(params), errRead.Error())
		return ""
	}
	return string(bytes)

}

func (agent *HttpAgent) Request(method string, path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error) {
	switch method {
	case http.MethodGet:
		response, err = agent.Get(path, header, timeoutMs, params)
		return
	case http.MethodPost:
		response, err = agent.Post(path, header, timeoutMs, params)
		return
	case http.MethodPut:
		response, err = agent.Put(path, header, timeoutMs, params)
		return
	case http.MethodDelete:
		response, err = agent.Delete(path, header, timeoutMs, params)
		return
	default:
		err = errors.New("not avaliable method")
		log.Printf("[ERROR]:request method[%s], path[%s],header:[%s],params:[%s], not avaliable method ", method, path, utils.ToJsonString(header), utils.ToJsonString(params))
	}
	return
}
func (agent *HttpAgent) Post(path string, header http.Header, timeoutMs uint64,
	params map[string]string) (response *http.Response, err error) {
	return post(path, header, timeoutMs, params)
}
func (agent *HttpAgent) Delete(path string, header http.Header, timeoutMs uint64,
	params map[string]string) (response *http.Response, err error) {
	return delete(path, header, timeoutMs, params)
}
func (agent *HttpAgent) Put(path string, header http.Header, timeoutMs uint64,
	params map[string]string) (response *http.Response, err error) {
	return put(path, header, timeoutMs, params)
}
