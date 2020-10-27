package version

import (
	"runtime"
	"time"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683670Test)
}

func _1599190683670Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		list3 := []models.SysDept{
			{DeptId: 1, ParentId: 0, DeptPath: "/0/1", DeptName: "爱拓科技", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 7, ParentId: 1, DeptPath: "/0/1/7", DeptName: "研发部", Sort: 1, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 8, ParentId: 1, DeptPath: "/0/1/8", DeptName: "运维部", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 9, ParentId: 1, DeptPath: "/0/1/9", DeptName: "客服部", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			{DeptId: 10, ParentId: 1, DeptPath: "/0/1/10", DeptName: "人力资源", Sort: 3, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "1", CreateBy: "1", UpdateBy: "1", BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		}

		list4 := []system.SysConfig{
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

		err := tx.Create(list3).Error
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
