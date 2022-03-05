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
	// 获取存储空间。
	input := &obs.PutFileInput{}
	input.Bucket = e.BucketName
	input.Key = yourObjectName
	input.SourceFile = localFile.(string)
	output, err := e.Client.(*obs.ObsClient).PutFile(input)

	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
		fmt.Printf("ETag:%s, StorageClass:%s\n", output.ETag, output.StorageClass)
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

func (e *HuaWeiOBS) GetTempToken() (string, error) {
	return "", nil
}
