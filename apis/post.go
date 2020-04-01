package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/utils"
	"net/http"
	"strconv"
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
	data.PostId, _ = utils.StringToInt64(id)

	data.PostName = c.Request.FormValue("postName")
	data.DataScope = utils.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.AssertErr(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
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
	Post.PostId, _ = strconv.ParseInt(c.Param("postId"), 10, 64)
	result, err := Post.Get()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
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
	pkg.AssertErr(err, "", 500)
	result, err := data.Create()
	pkg.AssertErr(err, "", -1)
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
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
	pkg.AssertErr(err, "", -1)
	result, err := data.Update(data.PostId)
	pkg.AssertErr(err, "", -1)
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
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
	id, err := utils.StringToInt64(c.Param("postId"))
	data.UpdateBy = utils.GetUserIdStr(c)
	_, err = data.Delete(id)
	pkg.AssertErr(err, "删除失败", 500)
	var res models.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}
