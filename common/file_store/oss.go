package file_store

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

type ALiYunOSS struct {
	Client          interface{}
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

func (e *ALiYunOSS) Setup() error {
	// TODO: 如果需要使用阿里云OSS请在此处填写对应信息
	e.AccessKeyId = ""
	e.AccessKeySecret = ""
	e.Endpoint = ""
	e.BucketName = ""

	client, err := oss.New(e.Endpoint, e.AccessKeyId, e.AccessKeySecret)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	e.Client = client
	return nil
}

// UpLoad 文件上传
func (e *ALiYunOSS) UpLoad(yourObjectName string, localFile string) error {
	err := e.Setup()
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// 获取存储空间。
	bucket, err := e.Client.(*oss.Client).Bucket(e.BucketName)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	// 设置分片大小为100 KB，指定分片上传并发数为3，并开启断点续传上传。
	// 其中<yourObjectName>与objectKey是同一概念，表示断点续传上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// "LocalFile"为filePath，100*1024为partSize。
	err = bucket.UploadFile(yourObjectName, localFile, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	return nil
}
