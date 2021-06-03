package apis

import (
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/google/uuid"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

type SysUser struct {
	api.Api
}

// GetSysUserList
// @Summary 列表用户信息数据
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "username"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysUser [get]
// @Security Bearer
func (e SysUser) GetSysUserList(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserSearch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysUser, 0)
	var count int64

	err = s.GetSysUserPage(&req, p, &list, &count)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetSysUser
// @Summary 获取用户
// @Description 获取JSON
// @Tags 用户
// @Param userId path int true "用户编码"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sysUser/{userId} [get]
// @Security Bearer
func (e SysUser) GetSysUser(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}
	var object models.SysUser
	//数据权限检查
	p := actions.GetPermissionFromContext(c)
	err = s.GetSysUser(&req, p, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}
	e.OK(object, "查看成功")
}

// InsertSysUser
// @Summary 创建用户
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysUserControl true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sysUser [post]
func (e SysUser) InsertSysUser(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	err = s.InsertSysUser(&req)
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// UpdateSysUser
// @Summary 修改用户数据
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysUserControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/sysuser/{userId} [put]
func (e SysUser) UpdateSysUser(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	req.SetUpdateBy(user.GetUserId(c))

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.UpdateSysUser(&req, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// DeleteSysUser
// @Summary 删除用户数据
// @Description 删除数据
// @Tags 用户
// @Param userId path int true "userId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/sysuser/{userId} [delete]
func (e SysUser) DeleteSysUser(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	// 设置编辑人
	req.SetUpdateBy(user.GetUserId(c))

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.RemoveSysUser(&req, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// InsetSysUserAvatar
// @Summary 修改头像
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/avatar [post]
func (e SysUser) InsetSysUserAvatar(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserControl{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}
	// 数据权限检查
	p := actions.GetPermissionFromContext(c)
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	guid := uuid.New().String()
	filPath := "static/uploadfile/" + guid + ".jpg"
	for _, file := range files {
		e.Logger.Debugf("upload avatar file: %s", file.Filename)
		// 上传文件至指定目录
		err = c.SaveUploadedFile(file, filPath)
		if err != nil {
			e.Logger.Errorf("save file error, %s", err.Error())
			e.Error(http.StatusInternalServerError, err, "")
			return
		}
	}
	req.UserId = p.UserId
	req.Avatar = "/" + filPath

	err = s.UpdateSysUser(&req, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(filPath, "修改成功")
}

// UpdateSysUserStatus 修改用户状态
// @Summary 修改用户状态
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.UpdateSysUserStatusReq true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/user/status [put]
func (e SysUser) UpdateSysUserStatus(c *gin.Context) {
	s := service.SysUser{}
	req := dto.UpdateSysUserStatusReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	req.SetUpdateBy(user.GetUserId(c))

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.UpdateSysUserStatus(&req, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// ResetSysUserPwd 重置用户密码
// @Summary 重置用户密码
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.ResetSysUserPwdReq true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/user/pwd/reset [put]
func (e SysUser) ResetSysUserPwd(c *gin.Context) {
	s := service.SysUser{}
	req := dto.ResetSysUserPwdReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	req.SetUpdateBy(user.GetUserId(c))

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.ResetSysUserPwd(&req, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// UpdateSysUserPwd
// @Summary 重置密码
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body dto.PassWord true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/user/pwd [post]
// @Security Bearer
func (e SysUser) UpdateSysUserPwd(c *gin.Context) {
	s := service.SysUser{}
	req := dto.PassWord{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.UpdateSysUserPwd(user.GetUserId(c), req.OldPassword, req.NewPassword, p)
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusForbidden, err, "密码修改失败")
		return
	}
	e.OK(nil, "密码修改成功")
}

// GetSysUserProfile
// @Summary 获取个人中心用户
// @Description 获取JSON
// @Tags 个人中心
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/user/profile [get]
// @Security Bearer
func (e SysUser) GetSysUserProfile(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserById{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	req.Id = user.GetUserId(c)

	user := new(models.SysUser)
	roles := make([]models.SysRole, 0)
	posts := make([]models.SysPost, 0)
	err = s.GetSysUserProfile(&req, user, &roles, &posts)
	if err != nil {
		e.Logger.Errorf("get user profile error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "获取用户信息失败")
		return
	}
	e.OK(gin.H{
		"user":  user,
		"roles": roles,
		"posts": posts,
	}, "查询成功")
}

func (e SysUser) GetInfo(c *gin.Context) {
	req := dto.SysUserById{}
	s := service.SysUser{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
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
		list, _ := RoleMenu.GetPermis(e.Orm)
		mp["permissions"] = list
		mp["buttons"] = list
	}

	var sysUser models.SysUser

	req.Id = user.GetUserId(c)

	err = s.GetSysUser(&req, p, &sysUser)
	if err != nil {
		e.Error(http.StatusUnauthorized, err, "登录失败")
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
	e.OK(mp, "")
}
