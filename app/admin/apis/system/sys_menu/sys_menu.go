package sys_menu

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
)

type SysMenu struct {
	apis.Api
}

// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menulist [get]
// @Security Bearer
func (e *SysMenu) GetSysMenuList(c *gin.Context) {
	log := e.GetLogger(c)
	d := new(dto.SysMenuSearch)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查询列表
	err = d.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	var list *[]system.SysMenu
	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Log = log
	serviceSysMenu.Orm = db
	list, err = serviceSysMenu.GetSysMenuPage(d)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, list, "查询成功")
}

// @Summary Menu详情数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menu/{id} [get]
// @Security Bearer
func (e *SysMenu) GetSysMenu(c *gin.Context) {
	control := new(dto.SysMenuById)

	log := e.GetLogger(c)

	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object system.SysMenu

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Log = log
	serviceSysMenu.Orm = db
	err = serviceSysMenu.GetSysMenu(control, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

// @Summary 创建菜单
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Param menuName formData string true "menuName"
// @Param Path formData string false "Path"
// @Param Action formData string true "Action"
// @Param Permission formData string true "Permission"
// @Param ParentId formData string true "ParentId"
// @Param IsDel formData string true "IsDel"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menu [post]
// @Security Bearer
func (e *SysMenu) InsertSysMenu(c *gin.Context) {
	control := new(dto.SysMenuControl)

	log := e.GetLogger(c)

	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//新增操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	// 设置创建人
	object.SetCreateBy(user.GetUserId(c))

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Orm = db
	serviceSysMenu.Log = log
	err = serviceSysMenu.InsertSysMenu(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

// @Summary 修改菜单
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Param id path int true "id"
// @Param data body dto.SysMenuControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/menu/{id} [put]
// @Security Bearer
func (e *SysMenu) UpdateSysMenu(c *gin.Context) {
	control := new(dto.SysMenuControl)

	log := e.GetLogger(c)

	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//更新操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.SetUpdateBy(user.GetUserId(c))

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Orm = db
	serviceSysMenu.Log = log
	err = serviceSysMenu.UpdateSysMenu(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

// @Summary 删除菜单
// @Description 删除数据
// @Tags 菜单
// @Param data body dto.SysMenuById true "body"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/menu/ [delete]
func (e *SysMenu) DeleteSysMenu(c *gin.Context) {
	control := new(dto.SysMenuById)

	log := e.GetLogger(c)

	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Errorf("Bind error: %s", err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Orm = db
	serviceSysMenu.Log = log
	err = serviceSysMenu.RemoveSysMenu(control)
	if err != nil {
		log.Errorf("RemoveSysMenu error, %s", err)
		e.Error(c, http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(c, control.GetId(), "删除成功")
}

// @Summary 根据角色名称获取菜单列表数据（左菜单使用）
// @Description 获取JSON
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menurole [get]
// @Security Bearer
func (e *SysMenu) GetMenuRole(c *gin.Context) {
	log := e.GetLogger(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Log = log
	serviceSysMenu.Orm = db
	result, err := serviceSysMenu.SetMenuRole(user.GetRoleName(c))

	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.OK(c, result, "")
}

// @Summary 获取角色对应的菜单id数组
// @Description 获取JSON
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menuids/{id} [get]
// @Security Bearer
func (e *SysMenu) GetMenuIDS(c *gin.Context) {
	log := e.GetLogger(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}
	var data models.RoleMenu
	data.RoleName = c.GetString("role")
	data.UpdateBy = user.GetUserIdStr(c)
	result, err := data.GetIDS(db)
	if err != nil {
		log.Errorf("GetIDS error, %s", err.Error())
		e.Error(c, http.StatusInternalServerError, err, "获取失败")
		return
	}
	e.OK(c, result, "")
}

//// GetMenuTreeRoleselect 角色修改中的菜单列表
//func (e *SysMenu) GetMenuTreeRoleselect(c *gin.Context) {
//	var Menu models.Menu
//	var SysRole models.SysRole
//
//	id, err := tools.StringToInt(c.Param("roleId"))
//	SysRole.RoleId = id
//	//var r *models.SysRole
//	r, err := SysRole.Get()
//
//	var result *[]models.MenuLable
//	menuIds := make([]int, 0)
//	if r.RoleKey != "admin" {
//		result, err = Menu.SetMenuLabel()
//		tools.HasError(err, "抱歉未找到相关信息", -1)
//		if id != 0 {
//			menuIds, err = SysRole.GetRoleMeunId()
//			tools.HasError(err, "抱歉未找到相关信息", -1)
//		}
//	}
//	app.Custum(c, gin.H{
//		"code":        200,
//		"menus":       result,
//		"checkedKeys": menuIds,
//	})
//}

// @Summary 获取菜单树
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menuTreeselect [get]
// @Security Bearer
func (e *SysMenu) GetMenuTreeSelect(c *gin.Context) {
	log := e.GetLogger(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	d := new(dto.SelectRole)

	err = c.BindUri(d)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Log = log
	serviceSysMenu.Orm = db
	result, err := serviceSysMenu.SetSysMenuLabel()
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}
	s := service.SysRole{}
	s.Log = log
	s.Orm = db
	menuIds, err := s.GetRoleMenuId(db, d.RoleId)
	if err != nil {
		log.Errorf("GetIDS error, %s", err.Error())
		e.Error(c, http.StatusInternalServerError, err, "")
		return
	}
	e.OK(c, gin.H{
		"menus":       result,
		"checkedKeys": menuIds,
	}, "获取成功")
}
