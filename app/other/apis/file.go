package apis

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"github.com/google/uuid"

	"go-admin/common/file_store"
)

type FileResponse struct {
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

const path = "static/uploadfile/"

type File struct {
	api.Api
}

// UploadFile 上传图片
// @Summary 上传图片
// @Description 获取JSON
// @Tags 公共接口
// @Accept multipart/form-data
// @Param type query string true "type" (1：单图，2：多图, 3：base64图片)
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/public/uploadFile [post]
// @Security Bearer
func (e File) UploadFile(c *gin.Context) {
	e.MakeContext(c)
	tag, _ := c.GetPostForm("type")
	urlPrefix := fmt.Sprintf("http://%s/", c.Request.Host)
	var fileResponse FileResponse

	switch tag {
	case "1": // 单图
		var done bool
		fileResponse, done = e.singleFile(c, fileResponse, urlPrefix)
		if done {
			return
		}
		e.OK(fileResponse, "上传成功")
		return
	case "2": // 多图
		multipartFile := e.multipleFile(c, urlPrefix)
		e.OK(multipartFile, "上传成功")
		return
	case "3": // base64
		fileResponse = e.baseImg(c, fileResponse, urlPrefix)
		e.OK(fileResponse, "上传成功")
	default:
		var done bool
		fileResponse, done = e.singleFile(c, fileResponse, urlPrefix)
		if done {
			return
		}
		e.OK(fileResponse, "上传成功")
		return
	}

}

func (e File) baseImg(c *gin.Context, fileResponse FileResponse, urlPerfix string) FileResponse {
	files, _ := c.GetPostForm("file")
	file2list := strings.Split(files, ",")
	ddd, _ := base64.StdEncoding.DecodeString(file2list[1])
	guid := uuid.New().String()
	fileName := guid + ".jpg"
	err := utils.IsNotExistMkDir(path)
	if err != nil {
		e.Error(500, errors.New(""), "初始化文件路径失败")
	}
	base64File := path + fileName
	_ = ioutil.WriteFile(base64File, ddd, 0666)
	typeStr := strings.Replace(strings.Replace(file2list[0], "data:", "", -1), ";base64", "", -1)
	fileResponse = FileResponse{
		Size:     pkg.GetFileSize(base64File),
		Path:     base64File,
		FullPath: urlPerfix + base64File,
		Name:     "",
		Type:     typeStr,
	}
	source, _ := c.GetPostForm("source")
	err = thirdUpload(source, fileName, base64File)
	if err != nil {
		e.Error(200, errors.New(""), "上传第三方失败")
		return fileResponse
	}
	if source != "1" {
		fileResponse.Path = "/static/uploadfile/" + fileName
		fileResponse.FullPath = "/static/uploadfile/" + fileName
	}
	return fileResponse
}

func (e File) multipleFile(c *gin.Context, urlPerfix string) []FileResponse {
	files := c.Request.MultipartForm.File["file"]
	source, _ := c.GetPostForm("source")
	var multipartFile []FileResponse
	for _, f := range files {
		guid := uuid.New().String()
		fileName := guid + utils.GetExt(f.Filename)

		err := utils.IsNotExistMkDir(path)
		if err != nil {
			e.Error(500, errors.New(""), "初始化文件路径失败")
		}
		multipartFileName := path + fileName
		err1 := c.SaveUploadedFile(f, multipartFileName)
		fileType, _ := utils.GetType(multipartFileName)
		if err1 == nil {
			err := thirdUpload(source, fileName, multipartFileName)
			if err != nil {
				e.Error(500, errors.New(""), "上传第三方失败")
			} else {
				fileResponse := FileResponse{
					Size:     pkg.GetFileSize(multipartFileName),
					Path:     multipartFileName,
					FullPath: urlPerfix + multipartFileName,
					Name:     f.Filename,
					Type:     fileType,
				}
				if source != "1" {
					fileResponse.Path = "/static/uploadfile/" + fileName
					fileResponse.FullPath = "/static/uploadfile/" + fileName
				}
				multipartFile = append(multipartFile, fileResponse)
			}
		}
	}
	return multipartFile
}

func (e File) singleFile(c *gin.Context, fileResponse FileResponse, urlPerfix string) (FileResponse, bool) {
	files, err := c.FormFile("file")

	if err != nil {
		e.Error(200, errors.New(""), "图片不能为空")
		return FileResponse{}, true
	}
	// 上传文件至指定目录
	guid := uuid.New().String()

	fileName := guid + utils.GetExt(files.Filename)

	err = utils.IsNotExistMkDir(path)
	if err != nil {
		e.Error(500, errors.New(""), "初始化文件路径失败")
	}
	singleFile := path + fileName
	_ = c.SaveUploadedFile(files, singleFile)
	fileType, _ := utils.GetType(singleFile)
	fileResponse = FileResponse{
		Size:     pkg.GetFileSize(singleFile),
		Path:     singleFile,
		FullPath: urlPerfix + singleFile,
		Name:     files.Filename,
		Type:     fileType,
	}
	//source, _ := c.GetPostForm("source")
	//err = thirdUpload(source, fileName, singleFile)
	//if err != nil {
	//	e.Error(200, errors.New(""), "上传第三方失败")
	//	return FileResponse{}, true
	//}
	fileResponse.Path = "/static/uploadfile/" + fileName
	fileResponse.FullPath = "/static/uploadfile/" + fileName
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