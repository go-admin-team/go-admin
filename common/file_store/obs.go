package file_store

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"log"
)

type HuaWeiOBS struct {
	Client     interface{}
	BucketName string
}

func (e *HuaWeiOBS) Setup(endpoint, accessKeyID, accessKeySecret, BucketName string, options ...ClientOption) error {
	// 创建ObsClient结构体
	client, err := obs.New(accessKeyID, accessKeySecret, endpoint)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	e.Client = client
	e.BucketName = BucketName
	return nil
}

// UpLoad 文件上传
// yourObjectName 文件路径名称，与objectKey是同一概念，表示断点续传上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg
func (e *HuaWeiOBS) UpLoad(yourObjectName string, localFile interface{}) error {
	client, ok := e.Client.(*obs.ObsClient)
	if !ok {
		return notConfigured(HuaweiOBS)
	}
	source, ok := localFile.(string)
	if !ok {
		return fmt.Errorf("obs upload wants a path, got %T", localFile)
	}

	// 获取存储空间。
	input := &obs.PutFileInput{}
	input.Bucket = e.BucketName
	input.Key = yourObjectName
	input.SourceFile = source
	output, err := client.PutFile(input)
	if err != nil {
		// The error used to be printed and nil returned, so a failed upload
		// reported success to the caller.
		if obsError, ok := err.(obs.ObsError); ok {
			return fmt.Errorf("obs upload: %s: %s", obsError.Code, obsError.Message)
		}
		return fmt.Errorf("obs upload: %w", err)
	}
	log.Printf("obs upload ok, requestId=%s etag=%s", output.RequestId, output.ETag)
	return nil
}

func (e *HuaWeiOBS) GetTempToken() (string, error) {
	return "", nil
}
