package vo

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-08 10:31
**/

type ConfigParam struct {
	DataId   string `param:"dataId"`
	Group    string `param:"group"`
	Content  string `param:"content"`
	Tag      string `param:"tag"`
	AppName  string `param:"appName"`
	OnChange func(namespace, group, dataId, data string)
}
