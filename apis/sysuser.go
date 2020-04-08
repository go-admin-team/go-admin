package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/utils"
	"log"
	"net/http"
	"strconv"
)

// @Summary 列表数据
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "username"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysUserList [get]
// @Security Bearer
func GetSysUserList(c *gin.Context) {
	var data models.SysUser
	var err error
	var pageSize = 10
	var pageIndex = 1

	size := c.Request.FormValue("pageSize")
	if size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	index := c.Request.FormValue("pageIndex")
	if index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.Username = c.Request.FormValue("userName")

	postId := c.Request.FormValue("postId")
	data.PostId, _ = strconv.ParseInt(postId, 10, 64)

	deptId := c.Request.FormValue("deptId")
	data.DeptId, _ = strconv.ParseInt(deptId, 10, 64)

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

// @Summary 获取用户
// @Description 获取JSON
// @Tags 用户
// @Param userId path int true "用户编码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sysUser/{userId} [get]
// @Security
func GetSysUser(c *gin.Context) {
	var SysUser models.SysUser
	SysUser.Id, _ = strconv.ParseInt(c.Param("userId"), 10, 64)
	result, err := SysUser.Get()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)
	var SysRole models.SysRole
	var Post models.Post
	roles, err := SysRole.GetList()
	posts, err := Post.GetList()

	postIds := make([]int64, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int64, 0)
	roleIds = append(roleIds, result.RoleId)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    result,
		"postIds": postIds,
		"roleIds": roleIds,
		"roles":   roles,
		"posts":   posts,
	})
}

// @Summary 获取当前登录用户
// @Description 获取JSON
// @Tags 个人中心
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/user/profile [get]
// @Security
func GetSysUserProfile(c *gin.Context) {
	var SysUser models.SysUser
	userId := utils.GetUserIdStr(c)
	SysUser.Id, _ = strconv.ParseInt(userId, 10, 64)
	result, err := SysUser.Get()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)
	var SysRole models.SysRole
	var Post models.Post
	var Dept models.Dept
	//获取角色列表
	roles, err := SysRole.GetList()
	//获取职位列表
	posts, err := Post.GetList()
	//获取部门列表
	Dept.Deptid = result.DeptId
	dept, err := Dept.Get()

	postIds := make([]int64, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int64, 0)
	roleIds = append(roleIds, result.RoleId)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    result,
		"postIds": postIds,
		"roleIds": roleIds,
		"roles":   roles,
		"posts":   posts,
		"dept":    dept,
	})
}

// @Summary 获取用户角色和职位
// @Description 获取JSON
// @Tags 用户
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sysUser [get]
// @Security
func GetSysUserInit(c *gin.Context) {
	var SysRole models.SysRole
	var Post models.Post
	roles, err := SysRole.GetList()
	posts, err := Post.GetList()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)
	mp := make(map[string]interface{}, 2)
	mp["roles"] = roles
	mp["posts"] = posts
	var res models.Response
	res.Data = mp
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 创建用户
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysUser true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sysUser [post]
func InsertSysUser(c *gin.Context) {
	var sysuser models.SysUser
	err := c.BindWith(&sysuser, binding.JSON)
	pkg.AssertErr(err, "非法数据格式", 500)

	sysuser.CreateTime = utils.GetCurrntTime()
	sysuser.UpdateTime = utils.GetCurrntTime()
	sysuser.CreateBy = utils.GetUserIdStr(c)
	sysuser.IsDel = 0
	id, err := sysuser.Insert()
	var res models.Response
	if err != nil {
		res.Msg = "添加失败"
		c.JSON(http.StatusOK, res.ReturnError(200))
		return
	}
	res.Data = id
	res.Msg = "添加成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 修改用户数据
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysUser true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/sysuser/{userId} [put]
func UpdateSysUser(c *gin.Context) {
	var data models.SysUser
	err := c.BindWith(&data, binding.JSON)
	pkg.AssertErr(err, "数据解析失败", -1)
	data.UpdateBy = utils.GetUserIdStr(c)
	result, err := data.Update(data.Id)
	var res models.Response
	if err != nil || result.Id == 0 {
		res.Msg = "修改失败"
		c.JSON(http.StatusOK, res.ReturnError(200))
		return
	}
	res.Msg = "修改成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 删除用户数据
// @Description 删除数据
// @Tags 用户
// @Param userId path int true "userId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/sysuser/{userId} [delete]
func DeleteSysUser(c *gin.Context) {
	var data models.SysUser
	data.UpdateBy = utils.GetUserIdStr(c)
	IDS := utils.IdsStrToIdsInt64Group("userId", c)
	result, err := data.BatchDelete(IDS)
	if err != nil || !result {
		var res models.Response
		res.Msg = "删除失败"
		c.JSON(http.StatusOK, res.ReturnError(501))
		return
	}
	var res models.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 修改头像
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/profileAvatar [post]
func InsetSysUserAvatar(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	guid := uuid.New().String()
	filPath := "static/uploadfile/" + guid + ".jpg"
	for _, file := range files {
		log.Println(file.Filename)
		// 上传文件至指定目录
		_ = c.SaveUploadedFile(file, filPath)
	}
	sysuser := models.SysUser{}
	sysuser.Id = utils.GetUserId(c)
	sysuser.Avatar = "/" + filPath
	sysuser.UpdateBy = utils.GetUserIdStr(c)
	sysuser.Update(sysuser.Id)

	c.JSON(http.StatusCreated, gin.H{
		"code": 200,
		"data": filPath,
	})

}

func SysUserUpdatePwd(c *gin.Context) {
	var pwd models.SysUserPwd
	err := c.Bind(&pwd)
	pkg.AssertErr(err, "数据解析失败", 500)
	sysuser := models.SysUser{}
	sysuser.Id = utils.GetUserId(c)
	sysuser.SetPwd(pwd)
	c.JSON(http.StatusCreated, gin.H{
		"code": 200,
		"data": "密码修改成功",
	})

}
