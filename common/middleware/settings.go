package middleware

type UrlInfo struct {
	Url    string
	Method string
}

// CasbinExclude casbin 排除的路由列表
var CasbinExclude = []UrlInfo{
	{Url: "/api/v1/dict/type-option-select", Method: "GET"},
	{Url: "/api/v1/dict-data/option-select", Method: "GET"},
	{Url: "/api/v1/deptTree", Method: "GET"},
	{Url: "/api/v1/db/tables/page", Method: "GET"},
	{Url: "/api/v1/db/columns/page", Method: "GET"},
	{Url: "/api/v1/gen/toproject/:tableId", Method: "GET"},
	{Url: "/api/v1/gen/todb/:tableId", Method: "GET"},
	{Url: "/api/v1/gen/tabletree", Method: "GET"},
	{Url: "/api/v1/gen/preview/:tableId", Method: "GET"},
	{Url: "/api/v1/gen/apitofile/:tableId", Method: "GET"},
	{Url: "/api/v1/getCaptcha", Method: "GET"},
	{Url: "/api/v1/getinfo", Method: "GET"},
	{Url: "/api/v1/menuTreeselect", Method: "GET"},
	{Url: "/api/v1/menurole", Method: "GET"},
	{Url: "/api/v1/menuids", Method: "GET"},
	{Url: "/api/v1/roleMenuTreeselect/:roleId", Method: "GET"},
	{Url: "/api/v1/roleDeptTreeselect/:roleId", Method: "GET"},
	{Url: "/api/v1/refresh_token", Method: "GET"},
	{Url: "/api/v1/configKey/:configKey", Method: "GET"},
	{Url: "/api/v1/app-config", Method: "GET"},
	{Url: "/api/v1/user/profile", Method: "GET"},
	{Url: "/info", Method: "GET"},
	{Url: "/api/v1/login", Method: "POST"},
	{Url: "/api/v1/logout", Method: "POST"},
	{Url: "/api/v1/user/avatar", Method: "POST"},
	{Url: "/api/v1/user/pwd", Method: "PUT"},
	{Url: "/api/v1/metrics", Method: "GET"},
	{Url: "/api/v1/health", Method: "GET"},
	{Url: "/", Method: "GET"},
	{Url: "/api/v1/server-monitor", Method: "GET"},
	{Url: "/api/v1/public/uploadFile", Method: "POST"},
}
