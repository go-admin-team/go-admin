package file_store

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Zone string

const (
	// HuaDong 华东
	HuaDong Zone = "HuaDong"
	// HuaBei 华北
	HuaBei Zone = "HuaBei"
	// HuaNan 华南
	HuaNan Zone = "HuaNan"
	// BeiMei 北美
	BeiMei Zone = "BeiMei"
	// XinJiaPo 新加坡
	XinJiaPo Zone = "XinJiaPo"
)

type QiNiuKODO struct {
	Client     interface{}
	BucketName string
	cfg        storage.Config
	options    []ClientOption
}

func (e *QiNiuKODO) getToken() string {
	putPolicy := storage.PutPolicy{
		Scope: e.BucketName,
	}
	if len(e.options) > 0 && e.options[0]["Expires"] != nil {
		putPolicy.Expires = e.options[0]["Expires"].(uint64)
	}
	upToken := putPolicy.UploadToken(e.Client.(*qbox.Mac))
	return upToken
}

//Setup 装载
//endpoint sss
func (e *QiNiuKODO) Setup(endpoint, accessKeyID, accessKeySecret, BucketName string, options ...ClientOption) error {

	mac := qbox.NewMac(accessKeyID, accessKeySecret)
	// 获取存储空间。
	cfg := storage.Config{}
	// 空间对应的机房
	e.setZoneORDefault(cfg, options...)
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	e.Client = mac
	e.BucketName = BucketName
	e.cfg = cfg
	e.options = options
	return nil
}

// setZoneORDefault 设置Zone或者默认华东
func (e *QiNiuKODO) setZoneORDefault(cfg storage.Config, options ...ClientOption) {
	if len(options) > 0 && options[0]["Zone"] != nil {
		if _, ok := options[0]["Zone"].(Zone); !ok {
			cfg.Zone = &storage.ZoneHuadong
		}
		switch options[0]["Zone"].(Zone) {
		case HuaDong:
			cfg.Zone = &storage.ZoneHuadong
		case HuaBei:
			cfg.Zone = &storage.ZoneHuabei
		case HuaNan:
			cfg.Zone = &storage.ZoneHuanan
		case BeiMei:
			cfg.Zone = &storage.ZoneBeimei
		case XinJiaPo:
			cfg.Zone = &storage.ZoneXinjiapo
		default:
			cfg.Zone = &storage.ZoneHuadong
		}
	}
}

// UpLoad 文件上传
func (e *QiNiuKODO) UpLoad(yourObjectName string, localFile interface{}) error {

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&e.cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, e.getToken(), yourObjectName, localFile.(string), &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(ret.Key, ret.Hash)
	return nil
}

func (e *QiNiuKODO) GetTempToken() (string, error) {
	token := e.getToken()
	return token, nil
}
