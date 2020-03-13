package http_agent

import "net/http"

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-10 11:07
**/

//go:generate mockgen -destination ../../mock/mock_http_agent_interface.go -package mock -source=./http_agent_interface.go

type IHttpAgent interface {
	Get(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	Post(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	Delete(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	Put(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	RequestOnlyResult(method string, path string, header http.Header, timeoutMs uint64, params map[string]string) string
	Request(method string, path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
}
