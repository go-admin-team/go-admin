package global

// Status values written to sys_opera_log.status.
//
// They live here rather than in app/admin/service/dto because
// common/middleware/logger.go writes the operation-log message and needs them.
// A package promised as a stable contract must not compile-depend on a
// business module: a fork that replaces or drops app/admin would otherwise
// stop compiling common/middleware, which is not something a contract package
// is allowed to do. See docs/contract.md.
const (
	OperaStatusEnabled  = "1"
	OperaStatusDisabled = "2"
)
