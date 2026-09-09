package models

import "go-admin/common/models"

type SysMenu struct {
	MenuId     int       `json:"menuId" gorm:"primaryKey;autoIncrement"`
	MenuName   string    `json:"menuName" gorm:"size:128;"`
	Title      string    `json:"title" gorm:"size:128;"`
	Icon       string    `json:"icon" gorm:"size:128;"`
	Path       string    `json:"path" gorm:"size:128;"`
	Paths      string    `json:"paths" gorm:"size:128;"`
	MenuType   string    `json:"menuType" gorm:"size:1;"`
	Action     string    `json:"action" gorm:"size:16;"`
	Permission string    `json:"permission" gorm:"size:255;"`
	ParentId   int       `json:"parentId" gorm:"size:11;"`
	NoCache    bool      `json:"noCache" gorm:"size:8;"`
	Breadcrumb string    `json:"breadcrumb" gorm:"size:255;"`
	Component  string    `json:"component" gorm:"size:255;"`
	Sort       int       `json:"sort" gorm:"size:4;"`
	Visible    string    `json:"visible" gorm:"size:1;"`
	IsFrame    string    `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	SysApi     []SysApi  `json:"sysApi" gorm:"many2many:sys_menu_api_rule"`
	Apis       []int     `json:"apis" gorm:"-"`
	DataScope  string    `json:"dataScope" gorm:"-"`
	Params     string    `json:"params" gorm:"-"`
	RoleId     int       `gorm:"-"`
	Children   []SysMenu `json:"children,omitempty" gorm:"-"`
	IsSelect   bool      `json:"is_select" gorm:"-"`
	// AppCode identifies which application's seed.SeedMenus call wrote this
	// row; empty for the host's own built-in menus. NOT NULL DEFAULT '' for
	// the same reason sys_migration.app_code is (see contract/models.Migration):
	// AutoMigrate adding this column to an existing table leaves every
	// pre-existing row reading back as "" rather than NULL.
	AppCode string `json:"appCode" gorm:"type:varchar(64);not null;default:'';index:idx_sys_menu_app_code;comment:AppCode"`
	// SeedCode is the raw seed.MenuSpec.Code this row was created from, kept
	// so seedMenuTree can ask "did I already write this node" without
	// relying on MenuName's PascalCase concatenation, which is not
	// injective (see design doc §1.6). Nullable, unlike AppCode: every row
	// seed.SeedMenus writes sets a real value, but every pre-existing row -
	// the host's own hand-placed menus, and every app-seeded row written
	// before this column existed - has none, and there is no way to
	// backfill one that means anything. NULL is what lets an unbounded
	// number of those coexist under the same app_code without tripping the
	// unique index below: the database never treats two NULLs as equal, so
	// only rows that do carry a real code participate in the uniqueness
	// check at all.
	// uk_sys_menu_app_seed_code_del is created by the migration, not from
	// this tag, and deliberately: it covers (app_code, seed_code,
	// deleted_at), and this struct cannot say so. A named uniqueIndex tag
	// puts every field carrying that name into one index, and deleted_at
	// comes from the shared ModelTime embed, which no single model can add a
	// tag to. Naming it here anyway declared a unique index on seed_code
	// alone under the same name - stricter than the real one, forbidding two
	// applications from both having a "dir" node - and AutoMigrate on this
	// model would have created that one first, after which the migration's
	// HasIndex guard finds the name taken and leaves the wrong index in
	// place.
	SeedCode *string `json:"seedCode" gorm:"size:64;comment:raw MenuSpec.Code, null for rows not written through SeedMenus"`
	models.ControlBy
	models.ModelTime
}

type SysMenuSlice []SysMenu

func (x SysMenuSlice) Len() int           { return len(x) }
func (x SysMenuSlice) Less(i, j int) bool { return x[i].Sort < x[j].Sort }
func (x SysMenuSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (*SysMenu) TableName() string {
	return "sys_menu"
}

func (e *SysMenu) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysMenu) GetId() interface{} {
	return e.MenuId
}
