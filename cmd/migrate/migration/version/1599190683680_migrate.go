package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683680Test)
}

func _1599190683680Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var err error

		list := []models.Menu{
			{MenuId: 496, MenuName: "Sources", Title: "资源管理", Icon: "network", Path: "/sources", Paths: "/0/496", MenuType: "M", Action: "无", Permission: "", ParentId: 0, NoCache: true, Breadcrumb: "", Component: "Layout", Sort: 3, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "1"},
			{MenuId: 497, MenuName: "File", Title: "文件管理", Icon: "documentation", Path: "file-manage", Paths: "/0/496/497", MenuType: "C", Action: "", Permission: "", ParentId: 496, NoCache: true, Breadcrumb: "", Component: "/fileManage/index", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "1"},
			{MenuId: 498, MenuName: "", Title: "内容管理", Icon: "pass", Path: "/content", Paths: "/0/498", MenuType: "M", Action: "无", Permission: "", ParentId: 0, NoCache: true, Breadcrumb: "", Component: "Layout", Sort: 4, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "1"},
			{MenuId: 499, MenuName: "SysCategory", Title: "分类", Icon: "pass", Path: "syscategory", Paths: "/0/498/499", MenuType: "C", Action: "无", Permission: "syscategory:syscategory:list", ParentId: 498, NoCache: true, Breadcrumb: "", Component: "/syscategory/index", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 500, MenuName: "", Title: "分页获取分类", Icon: "pass", Path: "", Paths: "/0/498/499/500", MenuType: "F", Action: "无", Permission: "syscategory:syscategory:query", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 501, MenuName: "", Title: "创建分类", Icon: "pass", Path: "", Paths: "/0/498/499/501", MenuType: "F", Action: "无", Permission: "syscategory:syscategory:add", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 502, MenuName: "", Title: "修改分类", Icon: "pass", Path: "", Paths: "/0/498/499/502", MenuType: "F", Action: "无", Permission: "syscategory:syscategory:edit", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 503, MenuName: "", Title: "删除分类", Icon: "pass", Path: "", Paths: "/0/498/499/503", MenuType: "F", Action: "无", Permission: "syscategory:syscategory:remove", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 504, MenuName: "Category", Title: "分类", Icon: "bug", Path: "category", Paths: "/0/63/504", MenuType: "M", Action: "无", Permission: "", ParentId: 63, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 505, MenuName: "", Title: "分页获取分类", Icon: "bug", Path: "/api/v1/syscategoryList", Paths: "/0/63/504/505", MenuType: "A", Action: "GET", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 506, MenuName: "", Title: "根据id获取分类", Icon: "bug", Path: "/api/v1/syscategory/:id", Paths: "/0/63/504/506", MenuType: "A", Action: "GET", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 507, MenuName: "", Title: "创建分类", Icon: "bug", Path: "/api/v1/syscategory", Paths: "/0/63/504/507", MenuType: "A", Action: "POST", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 508, MenuName: "", Title: "修改分类", Icon: "bug", Path: "/api/v1/syscategory", Paths: "/0/63/504/508", MenuType: "A", Action: "PUT", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 509, MenuName: "", Title: "删除分类", Icon: "bug", Path: "/api/v1/syscategory/:id", Paths: "/0/63/504/509", MenuType: "A", Action: "DELETE", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 511, MenuName: "SysContent", Title: "内容管理", Icon: "pass", Path: "syscontent", Paths: "/0/498/511", MenuType: "C", Action: "无", Permission: "syscontent:syscontent:list", ParentId: 498, NoCache: true, Breadcrumb: "", Component: "/syscontent/index", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 512, MenuName: "", Title: "分页获取内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/512", MenuType: "F", Action: "无", Permission: "syscontent:syscontent:query", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 513, MenuName: "", Title: "创建内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/513", MenuType: "F", Action: "无", Permission: "syscontent:syscontent:add", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 514, MenuName: "", Title: "修改内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/514", MenuType: "F", Action: "无", Permission: "syscontent:syscontent:edit", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 515, MenuName: "", Title: "删除内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/515", MenuType: "F", Action: "无", Permission: "syscontent:syscontent:remove", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 516, MenuName: "Content", Title: "内容管理", Icon: "bug", Path: "content", Paths: "/0/63/516", MenuType: "M", Action: "无", Permission: "", ParentId: 63, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 517, MenuName: "", Title: "分页获取内容管理", Icon: "bug", Path: "/api/v1/syscontentList", Paths: "/0/63/516/517", MenuType: "A", Action: "GET", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 518, MenuName: "", Title: "根据id获取内容管理", Icon: "bug", Path: "/api/v1/syscontent/:id", Paths: "/0/63/516/518", MenuType: "A", Action: "GET", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 519, MenuName: "", Title: "创建内容管理", Icon: "bug", Path: "/api/v1/syscontent", Paths: "/0/63/516/519", MenuType: "A", Action: "POST", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 520, MenuName: "", Title: "修改内容管理", Icon: "bug", Path: "/api/v1/syscontent", Paths: "/0/63/516/520", MenuType: "A", Action: "PUT", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
			{MenuId: 521, MenuName: "", Title: "删除内容管理", Icon: "bug", Path: "/api/v1/syscontent/:id", Paths: "/0/63/516/521", MenuType: "A", Action: "DELETE", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", CreateBy: "1", UpdateBy: "1", IsFrame: "0"},
		}

		err = tx.Create(list).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
