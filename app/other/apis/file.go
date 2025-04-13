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
	urlPrefix := fmt.Sprintf("%s://%s/", "http", c.Request.Host)

	switch tag {
	case "1": // 单图
		e.handleSingleFile(c, urlPrefix)
	case "2": // 多图
		e.handleMultipleFiles(c, urlPrefix)
	case "3": // base64
		e.handleBase64File(c, urlPrefix)
	default:
		e.handleSingleFile(c, urlPrefix)
	}
}

func (e File) handleSingleFile(c *gin.Context, urlPrefix string) {
	fileResponse, done := e.singleFile(c, FileResponse{}, urlPrefix)
	if done {
		return
	}
	e.OK(fileResponse, "上传成功")
}

func (e File) handleMultipleFiles(c *gin.Context, urlPrefix string) {
	multipartFile := e.multipleFile(c, urlPrefix)
	e.OK(multipartFile, "上传成功")
}

func (e File) handleBase64File(c *gin.Context, urlPrefix string) {
	fileResponse := e.baseImg(c, FileResponse{}, urlPrefix)
	e.OK(fileResponse, "上传成功")
}

func (e File) baseImg(c *gin.Context, fileResponse FileResponse, urlPrefix string) FileResponse {
	files, _ := c.GetPostForm("file")
	file2list := strings.Split(files, ",")
	decodedData, _ := base64.StdEncoding.DecodeString(file2list[1])
	fileName := uuid.New().String() + ".jpg"

	if err := utils.IsNotExistMkDir(path); err != nil {
		e.Error(500, errors.New(""), "初始化文件路径失败")
		return fileResponse
	}

	base64File := path + fileName
	_ = ioutil.WriteFile(base64File, decodedData, 0666)
	typeStr := strings.Replace(strings.Replace(file2list[0], "data:", "", -1), ";base64", "", -1)

	fileResponse = e.buildFileResponse(base64File, urlPrefix, "", typeStr)
	source, _ := c.GetPostForm("source")

	if err := thirdUpload(source, fileName, base64File); err != nil {
		e.Error(200, errors.New(""), "上传第三方失败")
		return fileResponse
	}

	if source != "1" {
		fileResponse.Path = "/static/uploadfile/" + fileName
		fileResponse.FullPath = "/static/uploadfile/" + fileName
	}
	return fileResponse
}

func (e File) multipleFile(c *gin.Context, urlPrefix string) []FileResponse {
	files := c.Request.MultipartForm.File["file"]
	source, _ := c.GetPostForm("source")
	var multipartFile []FileResponse

	for _, f := range files {
		fileName := uuid.New().String() + utils.GetExt(f.Filename)

		if err := utils.IsNotExistMkDir(path); err != nil {
			e.Error(500, errors.New(""), "初始化文件路径失败")
			continue
		}

		multipartFileName := path + fileName
		if err := c.SaveUploadedFile(f, multipartFileName); err != nil {
			continue
		}

		fileType, _ := utils.GetType(multipartFileName)
		if err := thirdUpload(source, fileName, multipartFileName); err != nil {
			e.Error(500, errors.New(""), "上传第三方失败")
			continue
		}

		fileResponse := e.buildFileResponse(multipartFileName, urlPrefix, f.Filename, fileType)
		if source != "1" {
			fileResponse.Path = "/static/uploadfile/" + fileName
			fileResponse.FullPath = "/static/uploadfile/" + fileName
		}
		multipartFile = append(multipartFile, fileResponse)
	}
	return multipartFile
}

func (e File) singleFile(c *gin.Context, fileResponse FileResponse, urlPrefix string) (FileResponse, bool) {
	files, err := c.FormFile("file")
	if err != nil {
		e.Error(200, errors.New(""), "图片不能为空")
		return FileResponse{}, true
	}

	fileName := uuid.New().String() + utils.GetExt(files.Filename)
	if err := utils.IsNotExistMkDir(path); err != nil {
		e.Error(500, errors.New(""), "初始化文件路径失败")
		return FileResponse{}, true
	}

	singleFile := path + fileName
	if err := c.SaveUploadedFile(files, singleFile); err != nil {
		e.Error(500, errors.New(""), "文件保存失败")
		return FileResponse{}, true
	}

	fileType, _ := utils.GetType(singleFile)
	fileResponse = e.buildFileResponse(singleFile, urlPrefix, files.Filename, fileType)
	fileResponse.Path = "/static/uploadfile/" + fileName
	fileResponse.FullPath = "/static/uploadfile/" + fileName
	return fileResponse, false
}

func (e File) buildFileResponse(filePath, urlPrefix, fileName, fileType string) FileResponse {
	return FileResponse{
		Size:     pkg.GetFileSize(filePath),
		Path:     filePath,
		FullPath: urlPrefix + filePath,
		Name:     fileName,
		Type:     fileType,
	}
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
