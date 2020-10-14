package version

import (
	"runtime"
	"time"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683670Test)
}

func _1599190683670Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		list1 := []models.RoleMenu{
			{RoleId: 1, MenuId: 2, RoleName: "admin"},
			{RoleId: 1, MenuId: 3, RoleName: "admin"},
			{RoleId: 1, MenuId: 43, RoleName: "admin"},
			{RoleId: 1, MenuId: 44, RoleName: "admin"},
			{RoleId: 1, MenuId: 45, RoleName: "admin"},
			{RoleId: 1, MenuId: 46, RoleName: "admin"},
			{RoleId: 1, MenuId: 51, RoleName: "admin"},
			{RoleId: 1, MenuId: 52, RoleName: "admin"},
			{RoleId: 1, MenuId: 56, RoleName: "admin"},
			{RoleId: 1, MenuId: 57, RoleName: "admin"},
			{RoleId: 1, MenuId: 58, RoleName: "admin"},
			{RoleId: 1, MenuId: 59, RoleName: "admin"},
			{RoleId: 1, MenuId: 60, RoleName: "admin"},
			{RoleId: 1, MenuId: 61, RoleName: "admin"},
			{RoleId: 1, MenuId: 62, RoleName: "admin"},
			{RoleId: 1, MenuId: 63, RoleName: "admin"},
			{RoleId: 1, MenuId: 64, RoleName: "admin"},
			{RoleId: 1, MenuId: 66, RoleName: "admin"},
			{RoleId: 1, MenuId: 67, RoleName: "admin"},
			{RoleId: 1, MenuId: 68, RoleName: "admin"},
			{RoleId: 1, MenuId: 69, RoleName: "admin"},
			{RoleId: 1, MenuId: 70, RoleName: "admin"},
			{RoleId: 1, MenuId: 71, RoleName: "admin"},
			{RoleId: 1, MenuId: 72, RoleName: "admin"},
			{RoleId: 1, MenuId: 73, RoleName: "admin"},
			{RoleId: 1, MenuId: 74, RoleName: "admin"},
			{RoleId: 1, MenuId: 75, RoleName: "admin"},
			{RoleId: 1, MenuId: 76, RoleName: "admin"},
			{RoleId: 1, MenuId: 77, RoleName: "admin"},
			{RoleId: 1, MenuId: 78, RoleName: "admin"},
			{RoleId: 1, MenuId: 79, RoleName: "admin"},
			{RoleId: 1, MenuId: 80, RoleName: "admin"},
			{RoleId: 1, MenuId: 81, RoleName: "admin"},
			{RoleId: 1, MenuId: 82, RoleName: "admin"},
			{RoleId: 1, MenuId: 83, RoleName: "admin"},
			{RoleId: 1, MenuId: 84, RoleName: "admin"},
			{RoleId: 1, MenuId: 85, RoleName: "admin"},
			{RoleId: 1, MenuId: 86, RoleName: "admin"},
			{RoleId: 1, MenuId: 87, RoleName: "admin"},
			{RoleId: 1, MenuId: 89, RoleName: "admin"},
			{RoleId: 1, MenuId: 90, RoleName: "admin"},
			{RoleId: 1, MenuId: 91, RoleName: "admin"},
			{RoleId: 1, MenuId: 92, RoleName: "admin"},
			{RoleId: 1, MenuId: 93, RoleName: "admin"},
			{RoleId: 1, MenuId: 94, RoleName: "admin"},
			{RoleId: 1, MenuId: 95, RoleName: "admin"},
			{RoleId: 1, MenuId: 96, RoleName: "admin"},
			{RoleId: 1, MenuId: 97, RoleName: "admin"},
			{RoleId: 1, MenuId: 103, RoleName: "admin"},
			{RoleId: 1, MenuId: 104, RoleName: "admin"},
			{RoleId: 1, MenuId: 105, RoleName: "admin"},
			{RoleId: 1, MenuId: 106, RoleName: "admin"},
			{RoleId: 1, MenuId: 107, RoleName: "admin"},
			{RoleId: 1, MenuId: 108, RoleName: "admin"},
			{RoleId: 1, MenuId: 109, RoleName: "admin"},
			{RoleId: 1, MenuId: 110, RoleName: "admin"},
			{RoleId: 1, MenuId: 111, RoleName: "admin"},
			{RoleId: 1, MenuId: 112, RoleName: "admin"},
			{RoleId: 1, MenuId: 113, RoleName: "admin"},
			{RoleId: 1, MenuId: 114, RoleName: "admin"},
			{RoleId: 1, MenuId: 115, RoleName: "admin"},
			{RoleId: 1, MenuId: 116, RoleName: "admin"},
			{RoleId: 1, MenuId: 117, RoleName: "admin"},
			{RoleId: 1, MenuId: 118, RoleName: "admin"},
			{RoleId: 1, MenuId: 119, RoleName: "admin"},
			{RoleId: 1, MenuId: 120, RoleName: "admin"},
			{RoleId: 1, MenuId: 121, RoleName: "admin"},
			{RoleId: 1, MenuId: 122, RoleName: "admin"},
			{RoleId: 1, MenuId: 123, RoleName: "admin"},
			{RoleId: 1, MenuId: 138, RoleName: "admin"},
			{RoleId: 1, MenuId: 142, RoleName: "admin"},
			{RoleId: 1, MenuId: 201, RoleName: "admin"},
			{RoleId: 1, MenuId: 202, RoleName: "admin"},
			{RoleId: 1, MenuId: 203, RoleName: "admin"},
			{RoleId: 1, MenuId: 204, RoleName: "admin"},
			{RoleId: 1, MenuId: 205, RoleName: "admin"},
			{RoleId: 1, MenuId: 206, RoleName: "admin"},
			{RoleId: 1, MenuId: 211, RoleName: "admin"},
			{RoleId: 1, MenuId: 212, RoleName: "admin"},
			{RoleId: 1, MenuId: 213, RoleName: "admin"},
			{RoleId: 1, MenuId: 214, RoleName: "admin"},
			{RoleId: 1, MenuId: 215, RoleName: "admin"},
			{RoleId: 1, MenuId: 216, RoleName: "admin"},
			{RoleId: 1, MenuId: 217, RoleName: "admin"},
			{RoleId: 1, MenuId: 220, RoleName: "admin"},
			{RoleId: 1, MenuId: 221, RoleName: "admin"},
			{RoleId: 1, MenuId: 222, RoleName: "admin"},
			{RoleId: 1, MenuId: 223, RoleName: "admin"},
			{RoleId: 1, MenuId: 224, RoleName: "admin"},
			{RoleId: 1, MenuId: 225, RoleName: "admin"},
			{RoleId: 1, MenuId: 226, RoleName: "admin"},
			{RoleId: 1, MenuId: 227, RoleName: "admin"},
			{RoleId: 1, MenuId: 228, RoleName: "admin"},
			{RoleId: 1, MenuId: 229, RoleName: "admin"},
			{RoleId: 1, MenuId: 230, RoleName: "admin"},
			{RoleId: 1, MenuId: 231, RoleName: "admin"},
			{RoleId: 1, MenuId: 232, RoleName: "admin"},
			{RoleId: 1, MenuId: 233, RoleName: "admin"},
			{RoleId: 1, MenuId: 234, RoleName: "admin"},
			{RoleId: 1, MenuId: 235, RoleName: "admin"},
			{RoleId: 1, MenuId: 236, RoleName: "admin"},
			{RoleId: 1, MenuId: 237, RoleName: "admin"},
			{RoleId: 1, MenuId: 238, RoleName: "admin"},
			{RoleId: 1, MenuId: 239, RoleName: "admin"},
			{RoleId: 1, MenuId: 240, RoleName: "admin"},
			{RoleId: 1, MenuId: 241, RoleName: "admin"},
			{RoleId: 1, MenuId: 242, RoleName: "admin"},
			{RoleId: 1, MenuId: 243, RoleName: "admin"},
			{RoleId: 1, MenuId: 244, RoleName: "admin"},
			{RoleId: 1, MenuId: 245, RoleName: "admin"},
			{RoleId: 1, MenuId: 246, RoleName: "admin"},
			{RoleId: 1, MenuId: 247, RoleName: "admin"},
			{RoleId: 1, MenuId: 248, RoleName: "admin"},
			{RoleId: 1, MenuId: 249, RoleName: "admin"},
			{RoleId: 1, MenuId: 250, RoleName: "admin"},
			{RoleId: 1, MenuId: 251, RoleName: "admin"},
			{RoleId: 1, MenuId: 252, RoleName: "admin"},
			{RoleId: 1, MenuId: 253, RoleName: "admin"},
			{RoleId: 1, MenuId: 254, RoleName: "admin"},
			{RoleId: 1, MenuId: 255, RoleName: "admin"},
			{RoleId: 1, MenuId: 256, RoleName: "admin"},
			{RoleId: 1, MenuId: 257, RoleName: "admin"},
			{RoleId: 1, MenuId: 258, RoleName: "admin"},
			{RoleId: 1, MenuId: 259, RoleName: "admin"},
			{RoleId: 1, MenuId: 260, RoleName: "admin"},
			{RoleId: 1, MenuId: 261, RoleName: "admin"},
			{RoleId: 1, MenuId: 262, RoleName: "admin"},
			{RoleId: 1, MenuId: 263, RoleName: "admin"},
			{RoleId: 1, MenuId: 264, RoleName: "admin"},
			{RoleId: 1, MenuId: 267, RoleName: "admin"},
			{RoleId: 1, MenuId: 269, RoleName: "admin"},
			{RoleId: 1, MenuId: 459, RoleName: "admin"},
			{RoleId: 1, MenuId: 460, RoleName: "admin"},
			{RoleId: 1, MenuId: 461, RoleName: "admin"},
			{RoleId: 1, MenuId: 462, RoleName: "admin"},
			{RoleId: 1, MenuId: 463, RoleName: "admin"},
			{RoleId: 1, MenuId: 464, RoleName: "admin"},
			{RoleId: 1, MenuId: 465, RoleName: "admin"},
			{RoleId: 1, MenuId: 466, RoleName: "admin"},
			{RoleId: 1, MenuId: 467, RoleName: "admin"},
			{RoleId: 1, MenuId: 468, RoleName: "admin"},
			{RoleId: 1, MenuId: 469, RoleName: "admin"},
			{RoleId: 1, MenuId: 470, RoleName: "admin"},
			{RoleId: 1, MenuId: 471, RoleName: "admin"},
			{RoleId: 1, MenuId: 473, RoleName: "admin"},
			{RoleId: 1, MenuId: 474, RoleName: "admin"},
			{RoleId: 1, MenuId: 475, RoleName: "admin"},
			{RoleId: 1, MenuId: 476, RoleName: "admin"},
			{RoleId: 1, MenuId: 477, RoleName: "admin"},
			{RoleId: 1, MenuId: 478, RoleName: "admin"},
			{RoleId: 1, MenuId: 479, RoleName: "admin"},
			{RoleId: 1, MenuId: 480, RoleName: "admin"},
			{RoleId: 1, MenuId: 481, RoleName: "admin"},
			{RoleId: 1, MenuId: 482, RoleName: "admin"},
			{RoleId: 1, MenuId: 483, RoleName: "admin"},
		}
		list2 := []models.CasbinRule{
			{PType: "p", V0: "admin", V1: "/api/v1/menulist", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/menu", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/databytype/", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/menu", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/menu/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysUserList", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysUser/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysUser/", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysUser", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysUser", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysUser/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/user/profile", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/rolelist", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/role/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/role", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/role", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/role/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/configList", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/config/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/config", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/config", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/config/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/menurole", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/roleMenuTreeselect/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/menuTreeselect", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/rolemenu", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/rolemenu", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/rolemenu/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/deptList", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dept/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dept", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/dept", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/dept/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/datalist", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/databytype/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data/", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/typelist", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/postlist", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/post/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/post", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/post", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/post/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/menu/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/menuids", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/loginloglist", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/loginlog/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/operloglist", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/getinfo", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/roledatascope", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/roleDeptTreeselect/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/deptTree", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/configKey/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/logout", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/user/avatar", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/user/pwd", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/typeoptionselect", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysjob", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysjob/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysjob", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysjob", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/sysjob", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/syssettingList", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/syssetting/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/syssetting", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/syssetting", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/syssetting/:id", V2: "DELETE"},
		}

		list3 := []models.SysDept{
			{DeptId: 1, ParentId: 0, DeptPath: "/0/1", DeptName: "爱拓科技", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 7, ParentId: 1, DeptPath: "/0/1/7", DeptName: "研发部", Sort: 1, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 8, ParentId: 1, DeptPath: "/0/1/8", DeptName: "运维部", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 9, ParentId: 1, DeptPath: "/0/1/9", DeptName: "客服部", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 10, ParentId: 1, DeptPath: "/0/1/10", DeptName: "人力资源", Sort: 3, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "1", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		}

		list4 := []models.SysConfig{
			{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(), DeletedAt: gorm.DeletedAt{}}, ControlBy: common.ControlBy{CreateBy: 1, UpdateBy: 1}, ConfigName: "主框架页-默认皮肤样式名称", ConfigKey: "sys_index_skinName", ConfigValue: "skin-blue", ConfigType: "Y", Remark: "蓝色 skin-blue、绿色 skin-green、紫色 skin-purple、红色 skin-red、黄色 skin-yellow"},
			{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now(), DeletedAt: gorm.DeletedAt{}}, ControlBy: common.ControlBy{CreateBy: 1, UpdateBy: 1}, ConfigName: "用户管理-账号初始密码", ConfigKey: "sys.user.initPassword", ConfigValue: "123456", ConfigType: "Y", Remark: "初始化密码 123456"},
			{Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now(), DeletedAt: gorm.DeletedAt{}}, ControlBy: common.ControlBy{CreateBy: 1, UpdateBy: 1}, ConfigName: "主框架页-侧边栏主题", ConfigKey: "sys_index_sideTheme", ConfigValue: "theme-dark", ConfigType: "Y", Remark: "深色主题theme-dark，浅色主题theme-light"},
		}

		list5 := []models.Post{
			{PostId: 1, PostName: "首席执行官", PostCode: "CEO", Sort: 0, Status: "0", Remark: "首席执行官", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{PostId: 2, PostName: "首席技术执行官", PostCode: "CTO", Sort: 2, Status: "0", Remark: "首席技术执行官", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{PostId: 3, PostName: "首席运营官", PostCode: "COO", Sort: 3, Status: "0", Remark: "测试工程师", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		}

		list6 := []models.SysRole{
			{1, "系统管理员", "0", "admin", 1, "", "1", "", "", true, "", models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "", []int{}, []int{}},
		}

		list7 := []models.DictType{
			{DictId: 1, DictName: "系统开关", DictType: "sys_normal_disable", Status: "0", CreateBy: "1", UpdateBy: "1", Remark: "系统开关列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 2, DictName: "用户性别", DictType: "sys_user_sex", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "用户性别列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 3, DictName: "菜单状态", DictType: "sys_show_hide", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "菜单状态列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 4, DictName: "系统是否", DictType: "sys_yes_no", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "系统是否列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 5, DictName: "任务状态", DictType: "sys_job_status", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "任务状态列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 6, DictName: "任务分组", DictType: "sys_job_group", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "任务分组列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 7, DictName: "通知类型", DictType: "sys_notice_type", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "通知类型列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 8, DictName: "系统状态", DictType: "sys_common_status", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "登录状态列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 9, DictName: "操作类型", DictType: "sys_oper_type", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "操作类型列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictId: 10, DictName: "通知状态", DictType: "sys_notice_status", Status: "0", CreateBy: "1", UpdateBy: "", Remark: "通知状态列表", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		}

		list8 := []models.SysUser{
			{models.SysUserId{1}, models.LoginM{models.UserName{"admin"}, models.PassWord{"$2a$10$cKFFTCzGOvaIHHJY2K45Zuwt8TD6oPzYi4s5MzYIBAWCLL6ZhouP2"}}, models.SysUserB{"zhangwj", "13818888888", 1, "", "", "0", "1@qq.com", 1, 1, "1", "1", "", "0", models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "", ""}},
		}

		list9 := []models.DictData{
			{DictCode: 1, DictSort: 0, DictLabel: "正常", DictValue: "0", DictType: "sys_normal_disable", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "系统正常", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 2, DictSort: 0, DictLabel: "停用", DictValue: "1", DictType: "sys_normal_disable", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "系统停用", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 3, DictSort: 0, DictLabel: "男", DictValue: "0", DictType: "sys_user_sex", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "性别男", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 4, DictSort: 0, DictLabel: "女", DictValue: "1", DictType: "sys_user_sex", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "性别女", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 5, DictSort: 0, DictLabel: "未知", DictValue: "2", DictType: "sys_user_sex", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "性别未知", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 6, DictSort: 0, DictLabel: "显示", DictValue: "0", DictType: "sys_show_hide", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "显示菜单", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 7, DictSort: 0, DictLabel: "隐藏", DictValue: "1", DictType: "sys_show_hide", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "隐藏菜单", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 8, DictSort: 0, DictLabel: "是", DictValue: "Y", DictType: "sys_yes_no", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "系统默认是", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 9, DictSort: 0, DictLabel: "否", DictValue: "N", DictType: "sys_yes_no", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "系统默认否", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 10, DictSort: 0, DictLabel: "正常", DictValue: "2", DictType: "sys_job_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "正常状态", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 11, DictSort: 0, DictLabel: "停用", DictValue: "1", DictType: "sys_job_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "停用状态", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 12, DictSort: 0, DictLabel: "默认", DictValue: "DEFAULT", DictType: "sys_job_group", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "默认分组", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 13, DictSort: 0, DictLabel: "系统", DictValue: "SYSTEM", DictType: "sys_job_group", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "系统分组", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 14, DictSort: 0, DictLabel: "通知", DictValue: "1", DictType: "sys_notice_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "通知", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 15, DictSort: 0, DictLabel: "公告", DictValue: "2", DictType: "sys_notice_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "公告", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 16, DictSort: 0, DictLabel: "正常", DictValue: "0", DictType: "sys_common_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "正常状态", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 17, DictSort: 0, DictLabel: "关闭", DictValue: "1", DictType: "sys_common_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "关闭状态", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 18, DictSort: 0, DictLabel: "新增", DictValue: "1", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "新增操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 19, DictSort: 0, DictLabel: "修改", DictValue: "2", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "修改操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 20, DictSort: 0, DictLabel: "删除", DictValue: "3", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "删除操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 21, DictSort: 0, DictLabel: "授权", DictValue: "4", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "授权操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 22, DictSort: 0, DictLabel: "导出", DictValue: "5", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "导出操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 23, DictSort: 0, DictLabel: "导入", DictValue: "6", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "导入操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 24, DictSort: 0, DictLabel: "强退", DictValue: "7", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "强退操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 25, DictSort: 0, DictLabel: "生成代码", DictValue: "8", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "生成操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 26, DictSort: 0, DictLabel: "清空数据", DictValue: "9", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "清空操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 27, DictSort: 0, DictLabel: "成功", DictValue: "0", DictType: "sys_notice_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "成功状态", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 28, DictSort: 0, DictLabel: "失败", DictValue: "1", DictType: "sys_notice_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "", Remark: "失败状态", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 29, DictSort: 0, DictLabel: "登录", DictValue: "10", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "1", Remark: "登录操作", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 30, DictSort: 0, DictLabel: "退出", DictValue: "11", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "1", Remark: "", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DictCode: 31, DictSort: 0, DictLabel: "获取验证码", DictValue: "12", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", CreateBy: "1", UpdateBy: "1", Remark: "获取验证码", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		}

		list10 := []models.SysSetting{
			{1, "go-admin管理系统", "https://gitee.com/mydearzwj/image/raw/master/img/go-admin.png", models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		}

		list11 := []models.SysJob{
			{1, "接口测试", "DEFAULT", 1, "0/5 * * * * ", "http://localhost:8000", "", 1, 1, 1, 0, "", "", models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}, ""},
			{2, "函数测试", "DEFAULT", 2, "0/5 * * * * ", "ExamplesOne", "参数", 1, 1, 1, 0, "", "", models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}, ""},
		}

		err := tx.Create(list1).Error
		if err != nil {
			return err
		}
		err = tx.Create(list2).Error
		if err != nil {
			return err
		}

		err = tx.Create(list3).Error
		if err != nil {
			return err
		}

		err = tx.Create(list4).Error
		if err != nil {
			return err
		}

		err = tx.Create(list5).Error
		if err != nil {
			return err
		}

		err = tx.Create(list6).Error
		if err != nil {
			return err
		}

		err = tx.Create(list7).Error
		if err != nil {
			return err
		}

		err = tx.Create(list8).Error
		if err != nil {
			return err
		}

		err = tx.Create(list9).Error
		if err != nil {
			return err
		}

		err = tx.Create(list10).Error
		if err != nil {
			return err
		}

		err = tx.Create(list11).Error
		if err != nil {
			return err
		}

		if err := models.InitDb(tx); err != nil {

		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
