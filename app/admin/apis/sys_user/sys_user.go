package sys_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/google/uuid"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
	common "go-admin/common/models"
)

type SysUser struct {
	apis.Api
}

// @Summary 列表用户信息数据
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "username"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysUser [get]
// @Security Bearer
func (e *SysUser) GetSysUserList(c *gin.Context) {
	log := e.GetLogger(c)
	d := new(dto.SysUserSearch)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	req := d.Generate()

	//查询列表
	err = req.Bind(c)
	if err != nil {
		log.Warnf("Bind error: %s", err.Error())
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]system.SysUser, 0)
	var count int64
	serviceStudent := service.SysUser{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetSysUserPage(req, p, &list, &count)
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.PageOK(c, list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// @Summary 获取用户
// @Description 获取JSON
// @Tags 用户
// @Param userId path int true "用户编码"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sysUser/{userId} [get]
// @Security Bearer
func (e *SysUser) GetSysUser(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.SysUserById)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查看详情
	req := control.Generate()
	err = req.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object system.SysUser

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysUser := service.SysUser{}
	serviceSysUser.Log = log
	serviceSysUser.Orm = db
	err = serviceSysUser.GetSysUser(req, p, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

// @Summary 创建用户
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysUserControl true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sysUser [post]
func (e *SysUser) InsertSysUser(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.SysUserControl)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//新增操作
	req := control.Generate()
	err = req.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object common.ActiveRecord
	object, err = req.GenerateM()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	// 设置创建人
	object.SetCreateBy(user.GetUserId(c))

	serviceSysUser := service.SysUser{}
	serviceSysUser.Orm = db
	serviceSysUser.Log = log
	err = serviceSysUser.InsertSysUser(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

// @Summary 修改用户数据
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysUserControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/sysuser/{userId} [put]
func (e *SysUser) UpdateSysUser(c *gin.Context) {
	control := new(dto.SysUserControl)

	log := e.GetLogger(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	req := control.Generate()
	//更新操作
	err = req.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object common.ActiveRecord
	object, err = req.GenerateM()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.SetUpdateBy(user.GetUserId(c))

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysUser := service.SysUser{}
	serviceSysUser.Orm = db
	serviceSysUser.Log = log
	err = serviceSysUser.UpdateSysUser(object, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

// @Summary 删除用户数据
// @Description 删除数据
// @Tags 用户
// @Param userId path int true "userId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/sysuser/{userId} [delete]
func (e *SysUser) DeleteSysUser(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.SysUserById)

	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//删除操作
	req := control.Generate()
	err = req.Bind(c)
	if err != nil {
		log.Errorf("Bind error: %s", err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object common.ActiveRecord
	object, err = req.GenerateM()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}

	// 设置编辑人
	object.SetUpdateBy(user.GetUserId(c))

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysUser := service.SysUser{}
	serviceSysUser.Orm = db
	serviceSysUser.Log = log
	err = serviceSysUser.RemoveSysUser(req, object, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "删除成功")
}

// @Summary 修改头像
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/avatar [post]
func (e *SysUser) InsetSysUserAvatar(c *gin.Context) {
	log := e.GetLogger(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	guid := uuid.New().String()
	filPath := "static/uploadfile/" + guid + ".jpg"
	for _, file := range files {
		log.Debugf("upload avatar file: %s", file.Filename)
		// 上传文件至指定目录
		err = c.SaveUploadedFile(file, filPath)
		if err != nil {
			log.Errorf("save file error, %s", err.Error())
			e.Error(c, http.StatusInternalServerError, err, "")
			return
		}
	}

	object := &system.SysUser{
		UserId: p.UserId,
		Avatar: "/" + filPath,
	}
	serviceSysUser := service.SysUser{}
	serviceSysUser.Orm = db
	serviceSysUser.Log = log
	err = serviceSysUser.UpdateSysUser(object, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, filPath, "修改成功")
}

// @Summary 重置密码
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.PassWord true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/user/pwd [post]
// @Security Bearer
func (e *SysUser) SysUserUpdatePwd(c *gin.Context) {
	log := e.GetLogger(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	var pwd dto.PassWord
	err = c.Bind(&pwd)
	if err != nil {
		log.Errorf("Bind error: %s", err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysUser := service.SysUser{}
	serviceSysUser.Orm = db
	serviceSysUser.Log = log
	err = serviceSysUser.UpdateSysUserPwd(user.GetUserId(c), pwd.OldPassword, pwd.NewPassword, p)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusForbidden, err, "密码修改失败")
		return
	}
	e.OK(c, nil, "密码修改成功")
}

// @Summary 获取个人中心用户
// @Description 获取JSON
// @Tags 个人中心
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/user/profile [get]
// @Security Bearer
func (e *SysUser) GetSysUserProfile(c *gin.Context) {
	log := e.GetLogger(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	id := user.GetUserId(c)
	serviceSysUser := service.SysUser{}
	serviceSysUser.Log = log
	serviceSysUser.Orm = db
	user := new(system.SysUser)
	roles := make([]system.SysRole, 0)
	posts := make([]system.SysPost, 0)
	err = serviceSysUser.GetSysUserProfile(id, user, &roles, &posts)
	if err != nil {
		log.Errorf("get user profile error, %s", err.Error())
		e.Error(c, http.StatusInternalServerError, err, "获取用户信息失败")
		return
	}
	e.OK(c, gin.H{
		"user":  user,
		"roles": roles,
		"posts": posts,
	}, "查询成功")
	//var SysUser models.SysUser
	//userId := tools.GetUserIdStr(c)
	//SysUser.UserId, _ = tools.StringToInt(userId)
	//result, err := SysUser.Get()
	//tools.HasError(err, "抱歉未找到相关信息", -1)
	//var SysRole models.SysRole
	//var Post models.Post
	//var Dept models.SysDepts
	////获取角色列表
	//roles, err := SysRole.GetList()
	////获取职位列表
	//posts, err := Post.GetList()
	////获取部门列表
	//Dept.DeptId = result.DeptId
	//dept, err := Dept.Get()
	//
	//postIds := make([]int, 0)
	//postIds = append(postIds, result.PostId)
	//
	//roleIds := make([]int, 0)
	//roleIds = append(roleIds, result.RoleId)
	//
	//app.Custum(c, gin.H{
	//	"code":    200,
	//	"data":    result,
	//	"postIds": postIds,
	//	"roleIds": roleIds,
	//	"roles":   roles,
	//	"posts":   posts,
	//	"dept":    dept,
	//})
}

func (e *SysUser) GetInfo(c *gin.Context) {
	log := e.GetLogger(c)

	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	var roles = make([]string, 1)
	roles[0] = user.GetRoleName(c)

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"

	var buttons = make([]string, 1)
	buttons[0] = "*:*:*"

	RoleMenu := models.RoleMenu{}
	RoleMenu.RoleId = user.GetRoleId(c)

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	if user.GetRoleName(c) == "admin" || user.GetRoleName(c) == "系统管理员" {
		mp["permissions"] = permissions
		mp["buttons"] = buttons
	} else {
		list, _ := RoleMenu.GetPermis(db)
		mp["permissions"] = list
		mp["buttons"] = list
	}

	var sysUser system.SysUser
	req := new(dto.SysUserById)
	req.Id = user.GetUserId(c)
	serviceSysUser := service.SysUser{}
	serviceSysUser.Log = log
	serviceSysUser.Orm = db
	err = serviceSysUser.GetSysUser(req, p, &sysUser)
	if err != nil {
		e.Error(c, http.StatusUnauthorized, err, "登录失败")
		return
	}

	mp["introduction"] = " am a super administrator"

	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if sysUser.Avatar != "" {
		mp["avatar"] = sysUser.Avatar
	}
	mp["userName"] = sysUser.NickName
	mp["userId"] = sysUser.UserId
	mp["deptId"] = sysUser.DeptId
	mp["name"] = sysUser.NickName
	e.OK(c, mp, "")
}
