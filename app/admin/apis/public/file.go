package public

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	imgType "github.com/shamsher31/goimgtype"

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
			files, err := c.FormFile("file")
			if err != nil {
				app.Error(c, 200, errors.New(""), "图片不能为空")
				return
			}
			// 上传文件至指定目录
			guid := uuid.New().String()

			singleFile := "static/uploadfile/" + guid + utils.GetExt(files.Filename)
			_ = c.SaveUploadedFile(files, singleFile)
			fileType, _ := imgType.Get(singleFile)
			fileResponse = FileResponse{
				Size:     tools.GetFileSize(singleFile),
				Path:     singleFile,
				FullPath: urlPerfix + singleFile,
				Name:     files.Filename,
				Type:     fileType,
			}
			app.OK(c, fileResponse, "上传成功")
			return
		case "2": // 多图
			files := c.Request.MultipartForm.File["file"]
			var multipartFile []FileResponse
			for _, f := range files {
				guid := uuid.New().String()
				multipartFileName := "static/uploadfile/" + guid + utils.GetExt(f.Filename)
				e := c.SaveUploadedFile(f, multipartFileName)
				fileType, _ := imgType.Get(multipartFileName)
				if e == nil {
					multipartFile = append(multipartFile, FileResponse{
						Size:     tools.GetFileSize(multipartFileName),
						Path:     multipartFileName,
						FullPath: urlPerfix + multipartFileName,
						Name:     f.Filename,
						Type:     fileType,
					})
				}
			}
			app.OK(c, multipartFile, "上传成功")
			return
		case "3": // base64
			files, _ := c.GetPostForm("file")
			file2list := strings.Split(files, ",")
			ddd, _ := base64.StdEncoding.DecodeString(file2list[1])
			guid := uuid.New().String()
			base64File := "static/uploadfile/" + guid + ".jpg"
			_ = ioutil.WriteFile(base64File, ddd, 0666)
			typeStr := strings.Replace(strings.Replace(file2list[0], "data:", "", -1), ";base64", "", -1)
			fileResponse = FileResponse{
				Size:     tools.GetFileSize(base64File),
				Path:     base64File,
				FullPath: urlPerfix + base64File,
				Name:     "",
				Type:     typeStr,
			}
			app.OK(c, fileResponse, "上传成功")
		}
	}
}
