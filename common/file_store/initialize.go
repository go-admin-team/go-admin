package file_store

import "fmt"

type OXS struct {
	// Endpoint 访问域名
	Endpoint string
	// AccessKeyID AK
	AccessKeyID string
	// AccessKeySecret AKS
	AccessKeySecret string
	// BucketName 桶名称
	BucketName string
}

// Setup 配置文件存储driver
//
// A failed Setup used to be printed and the store returned anyway, so the
// caller received one whose Client was nil - which panicked on first use. The
// error is returned instead.
func (e *OXS) Setup(driver DriverType, options ...ClientOption) (FileStoreType, error) {
	var fileStore FileStoreType
	switch driver {
	case AliYunOSS:
		fileStore = new(ALiYunOSS)
	case HuaweiOBS:
		fileStore = new(HuaWeiOBS)
	case QiNiuKodo:
		fileStore = new(QiNiuKODO)
	default:
		return nil, fmt.Errorf("unsupported file store driver %q", driver)
	}
	if err := fileStore.Setup(e.Endpoint, e.AccessKeyID, e.AccessKeySecret, e.BucketName, options...); err != nil {
		return nil, fmt.Errorf("file store %s: %w", driver, err)
	}
	return fileStore, nil
}
