package config

var ExtConfig Extend

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name
//
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	AMap      AMap // 这里配置对应配置文件的结构即可
	FileStore FileStore
}

type AMap struct {
	Key string
}

// FileStore 对象存储。上传接口的 source 参数决定走哪一家：2 是阿里云，3 是七牛。
// 没有填的那一家在被请求时返回明确错误，而不是上传到别处或者崩溃。
//
// common/file_store 里还实现了华为云 OBS，但上传接口没有对应的 source 取值，
// 所以这里也不为它提供配置。
type FileStore struct {
	AliYun ObjectStore
	QiNiu  ObjectStore
}

type ObjectStore struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// Configured reports whether enough was filled in to attempt a connection.
func (o ObjectStore) Configured() bool {
	return o.Endpoint != "" && o.AccessKeyID != "" && o.AccessKeySecret != "" && o.BucketName != ""
}
