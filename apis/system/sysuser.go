package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/app"
	"go-admin/pkg/utils"
	"log"
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
	data.PostId, _ = utils.StringToInt(postId)

	deptId := c.Request.FormValue("deptId")
	data.DeptId, _ = utils.StringToInt(deptId)

	data.DataScope = utils.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
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
	SysUser.UserId, _ = utils.StringToInt(c.Param("userId"))
	result, err := SysUser.Get()
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	var SysRole models.SysRole
	var Post models.Post
	roles, err := SysRole.GetList()
	posts, err := Post.GetList()

	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)
	app.Custum(c, gin.H{
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
	SysUser.UserId, _ = utils.StringToInt(userId)
	result, err := SysUser.Get()
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	var SysRole models.SysRole
	var Post models.Post
	var Dept models.Dept
	//获取角色列表
	roles, err := SysRole.GetList()
	//获取职位列表
	posts, err := Post.GetList()
	//获取部门列表
	Dept.DeptId = result.DeptId
	dept, err := Dept.Get()

	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)

	app.Custum(c, gin.H{
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
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	mp := make(map[string]interface{}, 2)
	mp["roles"] = roles
	mp["posts"] = posts
	app.OK(c, mp, "")
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
	pkg.HasError(err, "非法数据格式", 500)

	sysuser.CreateBy = utils.GetUserIdStr(c)
	id, err := sysuser.Insert()
	pkg.HasError(err, "添加失败", 500)
	app.OK(c, id, "添加成功")
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
	pkg.HasError(err, "数据解析失败", -1)
	data.UpdateBy = utils.GetUserIdStr(c)
	result, err := data.Update(data.UserId)
	pkg.HasError(err, "修改失败", 500)
	app.OK(c, result, "修改成功")
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
	IDS := utils.IdsStrToIdsIntGroup("userId", c)
	result, err := data.BatchDelete(IDS)
	pkg.HasError(err, "删除失败", 500)
	app.OK(c, result, "删除成功")
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
	sysuser.UserId = utils.GetUserId(c)
	sysuser.Avatar = "/" + filPath
	sysuser.UpdateBy = utils.GetUserIdStr(c)
	sysuser.Update(sysuser.UserId)
	app.OK(c, filPath, "修改成功")
}

func SysUserUpdatePwd(c *gin.Context) {
	var pwd models.SysUserPwd
	err := c.Bind(&pwd)
	pkg.HasError(err, "数据解析失败", 500)
	sysuser := models.SysUser{}
	sysuser.UserId = utils.GetUserId(c)
	sysuser.SetPwd(pwd)
	app.OK(c, "", "密码修改成功")
}
