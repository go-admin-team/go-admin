package constant

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-07 15:13
**/

type ServerConfig struct {
	ContextPath string
	IpAddr      string
	Port        uint64
}

type ClientConfig struct {
	TimeoutMs            uint64
	ListenInterval       uint64
	BeatInterval         int64
	NamespaceId          string
	Endpoint             string
	AccessKey            string
	SecretKey            string
	CacheDir             string
	LogDir               string
	UpdateThreadNum      int
	NotLoadCacheAtStart  bool
	UpdateCacheWhenEmpty bool
	OpenKMS              bool
	RegionId             string
}
