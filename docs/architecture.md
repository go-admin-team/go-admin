# 架构说明

> 本文记录**为什么这样设计**与不易从代码直接读出的语义。
> 编码规范见根目录 `AGENTS.md`，标准写法见 `app/demo/`。

## 数据权限（DataScope）

`PermissionAction()` 中间件从数据库取出当前用户的 DataScope 存入 `gin.Context`；
Service 在查询时通过 `actions.Permission(tableName, p)` 这个 GORM Scope 追加 WHERE。

实现见 `common/actions/permission.go:62-79`，五档语义：

| 值 | 含义 | 过滤方式 |
|---|---|---|
| `1` 或其他 | 全部数据 | 不追加条件 |
| `2` | 本角色关联部门的数据 | `create_by` 属于 `sys_role_dept` 关联部门下的用户 |
| `3` | 本部门数据 | `create_by` 属于同部门用户 |
| `4` | 本部门及子部门 | 按 `sys_dept.dept_path` 前缀匹配 |
| `5` | 仅本人 | `create_by = 当前用户` |

**过滤依据是 `create_by` 字段**，因此参与数据权限的表必须内嵌 `models.ControlBy`。

总开关：`config/settings.yml` 的 `application.enabledp`。关闭时 `Permission` 直接返回原
查询，这也意味着**关闭开关后所有数据权限配置立即失效**，排查问题时先确认此项。

## 定时任务

两类任务，配置在 `sys_job` 表：

- **HTTP 任务** —— 按 Cron 表达式请求指定 URL
- **函数任务** —— 调用注册在 `app/jobs` 中的 Go 函数

自定义函数任务需实现 `JobExec` 接口（`app/jobs/type.go:10`）：

```go
type JobExec interface {
    Exec(arg interface{}) error
}
```

并注册进 `app/jobs/examples.go` 的 `jobList` 映射，键名与 `sys_job` 表中配置的调用目标
对应。

任务内需要数据库连接时，通过 `sdk.Runtime.GetDbByKey("*")` 获取，不要反向 import
`app/admin/service`。

## 配置扩展

业务自定义配置写在 `config/extend.go` 中的结构体，对应 `settings.yml` 的 `extend:`
节点，代码中通过 `config.ExtConfig.Xxx` 访问。

配置值支持环境变量占位：在 yml 中写 `${ENV_NAME}`，由 `go-admin-core` 的
`config/reader/preprocessor.go` 在加载时替换。敏感信息可借此避免写入文件。

## 多数据源

`WithContextDb` 中间件按请求解析出对应的数据库连接放入上下文，Service 通过 `e.Orm`
取用。**不要使用全局 DB 变量** —— 那会绕过多租户隔离。

多库配置见 `settings.yml` 的 `databases` 与 `registers` 节点，后者用于 dbresolver
读写分离。

## 数据库迁移

迁移文件放 `cmd/migrate/migration/version-local/`，文件名前 13 位为 Unix 毫秒时间戳，
框架按文件名升序执行，已执行版本记录在 `sys_migration` 表。

```bash
go run main.go migrate -c config/settings.yml -g   # 生成骨架
```

**已执行过的迁移文件不可修改** —— 版本号已入表，改动不会重跑。需要调整时新建迁移。

## 构建注意

- 默认构建禁用 CGO；使用 SQLite 需 `make build-sqlite`（带 `-tags sqlite3`）
- `mode: prod` 时不注册 Swagger 路由
- dev 模式下 JWT 超时被设为极大值，生产部署前务必确认 `mode` 与 `jwt.secret`

## 参考

- 编码规范：根目录 `AGENTS.md`
- 标准 CRUD 模块：`app/demo/`
- API 文档：`go generate` 生成到 `docs/admin/`，dev 模式下访问
  `/swagger/admin/index.html`
