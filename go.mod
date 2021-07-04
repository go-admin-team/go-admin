module go-admin

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alibaba/sentinel-golang v0.6.1
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/bytedance/go-tagexpr/v2 v2.7.12
	github.com/casbin/casbin/v2 v2.25.1
	github.com/gin-gonic/gin v1.7.2
	github.com/go-admin-team/go-admin-core v1.3.6
	github.com/go-admin-team/go-admin-core/sdk v1.3.6
	github.com/google/uuid v1.2.0
	github.com/mssola/user_agent v0.5.2
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.11.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil v3.21.5+incompatible
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	github.com/unrolled/secure v1.0.8
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	gorm.io/driver/mysql v1.0.4-0.20201206014609-ae5fd10184f6
	gorm.io/driver/postgres v1.0.6-0.20201208020313-1ed927cfab53
	gorm.io/driver/sqlite v1.1.5-0.20201206014648-c84401fbe3ba
	gorm.io/gorm v1.21.11
)

//replace (
//	github.com/go-admin-team/go-admin-core => ../go-admin-core
//	github.com/go-admin-team/go-admin-core/sdk => ../go-admin-core/sdk
//)
