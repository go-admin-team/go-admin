package models

//type Menu struct {
//	MenuId     int    `json:"menuId" gorm:"primaryKey;autoIncrement"`
//	MenuName   string `json:"menuName" gorm:"size:128;"`
//	Title      string `json:"title" gorm:"size:128;"`
//	Icon       string `json:"icon" gorm:"size:128;"`
//	Path       string `json:"path" gorm:"size:128;"`
//	Paths      string `json:"paths" gorm:"size:128;"`
//	MenuType   string `json:"menuType" gorm:"size:1;"`
//	Action     string `json:"action" gorm:"size:16;"`
//	Permission string `json:"permission" gorm:"size:255;"`
//	ParentId   int    `json:"parentId" gorm:"size:11;"`
//	NoCache    bool   `json:"noCache" gorm:"size:8;"`
//	Breadcrumb string `json:"breadcrumb" gorm:"size:255;"`
//	Component  string `json:"component" gorm:"size:255;"`
//	Sort       int    `json:"sort" gorm:"size:4;"`
//	Visible    string `json:"visible" gorm:"size:1;"`
//	CreateBy   string `json:"createBy" gorm:"size:128;"`
//	UpdateBy   string `json:"updateBy" gorm:"size:128;"`
//	IsFrame    string `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
//	DataScope  string `json:"dataScope" gorm:"-"`
//	Params     string `json:"params" gorm:"-"`
//	RoleId     int    `gorm:"-"`
//	Children   []Menu `json:"children" gorm:"-"`
//	IsSelect   bool   `json:"is_select" gorm:"-"`
//
//	models.ModelTime
//}
//
//func (Menu) TableName() string {
//	return "sys_menu"
//}
//
//type MenuLable struct {
//	Id       int         `json:"id" gorm:"-"`
//	Label    string      `json:"label" gorm:"-"`
//	Children []MenuLable `json:"children" gorm:"-"`
//}
//
//type Menus struct {
//	MenuId     int    `json:"menuId" gorm:"column:menu_id;primaryKey;autoIncrement;"`
//	MenuName   string `json:"menuName" gorm:"column:menu_name"`
//	Title      string `json:"title" gorm:"column:title"`
//	Icon       string `json:"icon" gorm:"column:icon"`
//	Path       string `json:"path" gorm:"column:path"`
//	MenuType   string `json:"menuType" gorm:"column:menu_type"`
//	Action     string `json:"action" gorm:"column:action"`
//	Permission string `json:"permission" gorm:"column:permission"`
//	ParentId   int    `json:"parentId" gorm:"column:parent_id"`
//	NoCache    bool   `json:"noCache" gorm:"column:no_cache"`
//	Breadcrumb string `json:"breadcrumb" gorm:"column:breadcrumb"`
//	Component  string `json:"component" gorm:"column:component"`
//	Sort       int    `json:"sort" gorm:"column:sort"`
//
//	Visible  string `json:"visible" gorm:"column:visible"`
//	Children []Menu `json:"children" gorm:"-"`
//
//	CreateBy  string `json:"createBy" gorm:"column:create_by"`
//	UpdateBy  string `json:"updateBy" gorm:"column:update_by"`
//	DataScope string `json:"dataScope" gorm:"-"`
//	Params    string `json:"params" gorm:"-"`
//	BaseModel
//}
//
//func (Menus) TableName() string {
//	return "sys_menu"
//}
//
//type MenuRole struct {
//	Menus
//	IsSelect bool `json:"is_select" gorm:"-"`
//}
//
//type MS []Menu
//
//func (e *Menu) Get(tx *gorm.DB) (Menus []Menu, err error) {
//	table := tx.Table(e.TableName())
//	if e.MenuName != "" {
//		table = table.Where("menu_name = ?", e.MenuName)
//	}
//	if e.Path != "" {
//		table = table.Where("path = ?", e.Path)
//	}
//	if e.Action != "" {
//		table = table.Where("action = ?", e.Action)
//	}
//	if e.MenuType != "" {
//		table = table.Where("menu_type = ?", e.MenuType)
//	}
//
//	if err = table.Order("sort").Find(&Menus).Error; err != nil {
//		return
//	}
//	return
//}
//
//func (e *Menu) GetPage(tx *gorm.DB) (Menus []Menu, err error) {
//	table := tx.Table(e.TableName())
//	if e.MenuName != "" {
//		table = table.Where("menu_name = ?", e.MenuName)
//	}
//	if e.Title != "" {
//		table = table.Where("title = ?", e.Title)
//	}
//	if e.Visible != "" {
//		table = table.Where("visible = ?", e.Visible)
//	}
//	if e.MenuType != "" {
//		table = table.Where("menu_type = ?", e.MenuType)
//	}
//
//	// 数据权限控制
//	dataPermission := new(DataPermission)
//	dataPermission.UserId, _ = pkg.StringToInt(e.DataScope)
//	table, err = dataPermission.GetDataScope("sys_menu", table)
//	if err != nil {
//		return nil, err
//	}
//	if err = table.Order("sort").Find(&Menus).Error; err != nil {
//		return
//	}
//	return
//}
//
//func (e *Menu) Create(tx *gorm.DB) (id int, err error) {
//	result := tx.Table(e.TableName()).Create(&e)
//	if result.Error != nil {
//		err = result.Error
//		return
//	}
//	err = InitPaths(tx, e)
//	if err != nil {
//		return
//	}
//	id = e.MenuId
//	return
//}
//
//func InitPaths(tx *gorm.DB, menu *Menu) (err error) {
//	parentMenu := new(Menu)
//	if menu.ParentId != 0 {
//		tx.Table("sys_menu").Where("menu_id = ?", menu.ParentId).First(parentMenu)
//		if parentMenu.Paths == "" {
//			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
//			return
//		}
//		menu.Paths = parentMenu.Paths + "/" + pkg.IntToString(menu.MenuId)
//	} else {
//		menu.Paths = "/0/" + pkg.IntToString(menu.MenuId)
//	}
//	tx.Table("sys_menu").Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths)
//	return
//}
