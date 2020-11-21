package public

import (
	"encoding/base64"
	"errors"
	"fmt"
	"go-admin/common/file_store"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	//imgType "github.com/shamsher31/goimgtype"

	"go-admin/pkg/utils"
	"go-admin/tools"
	"go-admin/tools/app"
)

type FileResponse struct {
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

// @Summary 上传图片
// @Description 获取JSON
// @Tags 公共接口
// @Accept multipart/form-data
// @Param type query string true "type" (1：单图，2：多图, 3：base64图片)
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/public/uploadFile [post]
func UploadFile(c *gin.Context) {
	tag, _ := c.GetPostForm("type")
	urlPerfix := fmt.Sprintf("http://%s/", c.Request.Host)
	var fileResponse FileResponse
	if tag == "" {
		app.Error(c, 200, errors.New(""), "缺少标识")
		return
	} else {
		switch tag {
		case "1": // 单图
			fileResponse, done := singleFile(c, fileResponse, urlPerfix)
			if done {
				return
			}
			app.OK(c, fileResponse, "上传成功")
			return
		case "2": // 多图
			multipartFile := multipleFile(c, urlPerfix)
			app.OK(c, multipartFile, "上传成功")
			return
		case "3": // base64
			fileResponse = baseImg(c, fileResponse, urlPerfix)
			app.OK(c, fileResponse, "上传成功")
		}
	}
}

func baseImg(c *gin.Context, fileResponse FileResponse, urlPerfix string) FileResponse {
	files, _ := c.GetPostForm("file")
	file2list := strings.Split(files, ",")
	ddd, _ := base64.StdEncoding.DecodeString(file2list[1])
	guid := uuid.New().String()
	fileName := guid + ".jpg"
	base64File := "static/uploadfile/" + fileName
	_ = ioutil.WriteFile(base64File, ddd, 0666)
	typeStr := strings.Replace(strings.Replace(file2list[0], "data:", "", -1), ";base64", "", -1)
	fileResponse = FileResponse{
		Size:     tools.GetFileSize(base64File),
		Path:     base64File,
		FullPath: urlPerfix + base64File,
		Name:     "",
		Type:     typeStr,
	}
	source, _ := c.GetPostForm("source")
	err := thirdUpload(source, fileName, base64File)
	if err != nil {
		app.Error(c, 200, errors.New(""), "上传第三方失败")
		return fileResponse
	}
	if source != "1" {
		fileResponse.Path = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
		fileResponse.FullPath = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
	}
	return fileResponse
}

func multipleFile(c *gin.Context, urlPerfix string) []FileResponse {
	files := c.Request.MultipartForm.File["file"]
	source, _ := c.GetPostForm("source")
	var multipartFile []FileResponse
	for _, f := range files {
		guid := uuid.New().String()
		fileName := guid + utils.GetExt(f.Filename)
		multipartFileName := "static/uploadfile/" + fileName
		e := c.SaveUploadedFile(f, multipartFileName)
		fileType, _ := utils.GetType(multipartFileName)
		if e == nil {
			err := thirdUpload(source, fileName, multipartFileName)
			if err != nil {
				app.Error(c, 200, errors.New(""), "上传第三方失败")
			} else {
				fileResponse := FileResponse{
					Size:     tools.GetFileSize(multipartFileName),
					Path:     multipartFileName,
					FullPath: urlPerfix + multipartFileName,
					Name:     f.Filename,
					Type:     fileType,
				}
				if source != "1" {
					fileResponse.Path = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
					fileResponse.FullPath = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
				}
				multipartFile = append(multipartFile, fileResponse)
			}
		}
	}
	return multipartFile
}

func singleFile(c *gin.Context, fileResponse FileResponse, urlPerfix string) (FileResponse, bool) {
	files, err := c.FormFile("file")

	if err != nil {
		app.Error(c, 200, errors.New(""), "图片不能为空")
		return FileResponse{}, true
	}
	// 上传文件至指定目录
	guid := uuid.New().String()

	fileName := guid + utils.GetExt(files.Filename)
	singleFile := "static/uploadfile/" + fileName
	_ = c.SaveUploadedFile(files, singleFile)
	fileType, _ := utils.GetType(singleFile)
	fileResponse = FileResponse{
		Size:     tools.GetFileSize(singleFile),
		Path:     singleFile,
		FullPath: urlPerfix + singleFile,
		Name:     files.Filename,
		Type:     fileType,
	}
	source, _ := c.GetPostForm("source")
	err = thirdUpload(source, fileName, singleFile)
	if err != nil {
		app.Error(c, 200, errors.New(""), "上传第三方失败")
		return FileResponse{}, true
	}
	fileResponse.Path = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
	fileResponse.FullPath = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
	return fileResponse, false
}

func thirdUpload(source string, name string, path string) error {
	switch source {
	case "2":
		return ossUpload("img/"+name, path)
	case "3":
		return qiniuUpload("img/"+name, path)
	}
	return nil
}

func ossUpload(name string, path string) error {
	oss := file_store.ALiYunOSS{}
	return oss.UpLoad(name, path)
}

func qiniuUpload(name string, path string) error {
	oss := file_store.ALiYunOSS{}
	return oss.UpLoad(name, path)
}
