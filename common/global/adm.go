package global

const (
	// Version go-admin version info
	Version = "2.4.0"
)

var (
	// Driver 数据库驱动
	//
	// Deprecated: common/dto.MakeCondition stopped reading this after PRD
	// 006 F2/F5 - it now takes the dialect from the *gorm.DB passed to the
	// scope it returns instead of this process-wide variable. Driver is
	// still set (common/database/initialize.go) and still readable for fork
	// code that reads it directly, but it is no longer this framework's own
	// path to the current SQL dialect.
	Driver string
)
