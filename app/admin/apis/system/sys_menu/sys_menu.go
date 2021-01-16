package sys_menu

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
	"go-admin/common/log"
	"go-admin/tools"
	"go-admin/tools/app"
	"net/http"
)

type SysMenu struct {
	apis.Api
}

func (e *SysMenu) GetSysMenuList(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	d := new(dto.SysMenuSearch)
	db, err := tools.GetOrm(c)
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

	var list *[]models.SysMenu
	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.MsgID = msgID
	serviceSysMenu.Orm = db
	list, err = serviceSysMenu.GetSysMenuPage(d)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, list, "查询成功")
}

func (e *SysMenu) GetSysMenu(c *gin.Context) {
	control := new(dto.SysMenuById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object models.SysMenu

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.MsgID = msgID
	serviceSysMenu.Orm = db
	err = serviceSysMenu.GetSysMenu(control, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

func (e *SysMenu) InsertSysMenu(c *gin.Context) {
	control := new(dto.SysMenuControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
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
	object.SetCreateBy(tools.GetUserId(c))

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Orm = db
	serviceSysMenu.MsgID = msgID
	err = serviceSysMenu.InsertSysMenu(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

func (e *SysMenu) UpdateSysMenu(c *gin.Context) {
	control := new(dto.SysMenuControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
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
	object.SetUpdateBy(tools.GetUserId(c))

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Orm = db
	serviceSysMenu.MsgID = msgID
	err = serviceSysMenu.UpdateSysMenu(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

func (e *SysMenu) DeleteSysMenu(c *gin.Context) {
	control := new(dto.SysMenuById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Errorf("MsgID[%s] Bind error: %s", msgID, err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.Orm = db
	serviceSysMenu.MsgID = msgID
	err = serviceSysMenu.RemoveSysMenu(control)
	if err != nil {
		log.Error(err)
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
	msgID := tools.GenerateMsgIDFromContext(c)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.MsgID = msgID
	serviceSysMenu.Orm = db
	result, err := serviceSysMenu.SetMenuRole(tools.GetRoleName(c))

	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	app.OK(c, result, "")
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
	var data models.RoleMenu
	data.RoleName = c.GetString("role")
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.GetIDS()
	tools.HasError(err, "获取失败", 500)
	app.OK(c, result, "")
}

// GetMenuTreeRoleselect 角色修改中的菜单列表
func (e *SysMenu) GetMenuTreeRoleselect(c *gin.Context) {
	var Menu models.Menu
	var SysRole models.SysRole

	id, err := tools.StringToInt(c.Param("roleId"))
	SysRole.RoleId = id
	//var r *models.SysRole
	r, err := SysRole.Get()

	var result *[]models.MenuLable
	menuIds := make([]int, 0)
	if r.RoleKey != "admin" {
		result, err = Menu.SetMenuLabel()
		tools.HasError(err, "抱歉未找到相关信息", -1)
		if id != 0 {
			menuIds, err = SysRole.GetRoleMeunId()
			tools.HasError(err, "抱歉未找到相关信息", -1)
		}
	}
	app.Custum(c, gin.H{
		"code":        200,
		"menus":       result,
		"checkedKeys": menuIds,
	})
}

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
	msgID := tools.GenerateMsgIDFromContext(c)
	d := new(dto.SysMenuSearch)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	err = d.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysMenu := service.SysMenu{}
	serviceSysMenu.MsgID = msgID
	serviceSysMenu.Orm = db
	result, err := serviceSysMenu.SetSysMenuLabel(d)

	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}
	e.OK(c, result, "")
}
