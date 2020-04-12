package system

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/app"
	"go-admin/pkg/utils"
)

// @Summary 职位列表数据
// @Description 获取JSON
// @Tags 职位
// @Param name query string false "name"
// @Param id query string false "id"
// @Param position query string false "position"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/post [get]
// @Security
func GetPostList(c *gin.Context) {
	var data models.Post
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.PostName = c.Request.FormValue("postName")
	id := c.Request.FormValue("postId")
	data.PostId, _ = utils.StringToInt(id)

	data.PostName = c.Request.FormValue("postName")
	data.DataScope = utils.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.HasError(err, "", -1)
	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

// @Summary 获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param postId path int true "postId"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/post/{postId} [get]
// @Security
func GetPost(c *gin.Context) {
	var Post models.Post
	Post.PostId, _ = utils.StringToInt(c.Param("postId"))
	result, err := Post.Get()
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c,result,"")
}

// @Summary 添加职位
// @Description 获取JSON
// @Tags 职位
// @Accept  application/json
// @Product application/json
// @Param data body models.Post true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/post [post]
// @Security Bearer
func InsertPost(c *gin.Context) {
	var data models.Post
	err := c.Bind(&data)
	data.CreateBy = utils.GetUserIdStr(c)
	pkg.HasError(err, "", 500)
	result, err := data.Create()
	pkg.HasError(err, "", -1)
	app.OK(c,result,"")
}

// @Summary 修改职位
// @Description 获取JSON
// @Tags 职位
// @Accept  application/json
// @Product application/json
// @Param data body models.Dept true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/post/ [put]
// @Security Bearer
func UpdatePost(c *gin.Context) {
	var data models.Post

	err := c.Bind(&data)
	data.UpdateBy = utils.GetUserIdStr(c)
	pkg.HasError(err, "", -1)
	result, err := data.Update(data.PostId)
	pkg.HasError(err, "", -1)
	app.OK(c,result,"修改成功")
}

// @Summary 删除职位
// @Description 删除数据
// @Tags 职位
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/post/{postId} [delete]
func DeletePost(c *gin.Context) {
	var data models.Post
	id, err := utils.StringToInt(c.Param("postId"))
	data.UpdateBy = utils.GetUserIdStr(c)
	_, err = data.Delete(id)
	pkg.HasError(err, "删除失败", 500)
	app.OK(c,"","删除成功")
}
