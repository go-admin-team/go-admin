# AGENTS.md — go-admin 后端

> 给 AI 编码工具与新贡献者的约定。**只写"不遵守就会出错"的规则**；技术栈版本以
> `go.mod` 为准，命令以 `Makefile` 为准，此处不复述，避免与代码脱节。
>
> 标准 CRUD 模块的完整写法见 **`app/demo/`** —— 那是可编译、有测试、CI 会跑的参照物。
> 本文与它冲突时，以 `app/demo/` 为准。

## 分层

```
Router  →  Api      →  Service      →  Model
路由注册    参数绑定      业务逻辑         GORM 结构体
中间件链    调用 Service  操作数据库       TableName()
```

对应目录：`app/{模块}/router|apis|service|models`，DTO 位于 `service/dto`。

**不可跨层**：Api 不直接操作 `Orm`，Service 不接触 `gin.Context`。

## 优先使用通用 Action

单表 CRUD **不要手写 Handler 与 Service**。`common/actions` 提供的五个
Action 已覆盖参数绑定、数据权限过滤、操作人注入、分页与错误响应：

```go
r := v1.Group("/demo-product").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
{
    m := &models.DemoProduct{}
    r.GET("",     actions.PermissionAction(), actions.IndexAction(m, new(dto.DemoProductSearch), func() interface{} {
        list := make([]models.DemoProduct, 0); return &list
    }))
    r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(dto.DemoProductById), func() interface{} {
        return &models.DemoProduct{}
    }))
    r.POST("",       actions.CreateAction(new(dto.DemoProductControl)))
    r.PUT("/:id",    actions.PermissionAction(), actions.UpdateAction(new(dto.DemoProductControl)))
    r.DELETE("",     actions.PermissionAction(), actions.DeleteAction(new(dto.DemoProductById)))
}
```

这样一个模块只需 **model + dto + router** 三个文件，完整示例见 `app/demo/`。

使用通用 Action 的前提：

- Model 实现 `models.ActiveRecord`（`Generate` / `GetId` / `TableName`）
- 列表 DTO 实现 `dto.Index`，增改删 DTO 实现 `dto.Control`
- **所有 `Generate()` 必须返回副本** —— Action 在并发请求间复用实例，
  就地返回会串数据（`app/demo` 的测试锁定了这一点）
- 详情/删除 DTO 内嵌 `dto.ObjectById` 即可继承 `Bind` 与 `GetId`，无需重写

仅当业务超出单表 CRUD（跨表事务、外部调用、复杂校验）时才自行编写 Handler
与 Service，写法见下。

## Api 层（仅在通用 Action 不适用时）

结构体嵌入 `api.Api`，链式初始化后**必须检查 `Errors`**：

```go
func (e SysPost) GetPage(c *gin.Context) {
    s := service.SysPost{}
    req := dto.SysPostPageReq{}
    err := e.MakeContext(c).MakeOrm().Bind(&req, binding.Form).MakeService(&s.Service).Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
    // ... 调用 s.GetPage(...)
    e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}
```

响应一律走 `e.OK` / `e.PageOK` / `e.Error`，不要自行 `c.JSON`。

## Service 层（仅在通用 Action 不适用时）

结构体嵌入 `service.Service`（持有 `Orm` 与 `Log`）。查询通过 Scopes 组合：

```go
err = e.Orm.Model(&data).Scopes(
    cDto.MakeCondition(c.GetNeedSearch()),          // 由 search tag 生成 WHERE
    cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
    actions.Permission(data.TableName(), p),        // 数据权限，列表/详情必须带
).Find(list).Limit(-1).Offset(-1).Count(count).Error
```

**遗漏 `actions.Permission` 会使数据权限配置静默失效** —— 这是最容易出的错。

错误一律 `return err` 向上传递，日志用 `e.Log.Errorf`，不使用 `panic`。

## DTO

搜索条件由 tag 声明，`MakeCondition` 据此拼 SQL：

```go
type SysPostPageReq struct {
    dto.Pagination `search:"-"`
    PostName string `form:"postName" search:"type:contains;column:post_name;table:sys_post"`
}
func (m *SysPostPageReq) GetNeedSearch() interface{} { return *m }
```

`type` 可选：`exact` `iexact` `contains` `gt` `gte` `lt` `lte` `order` `left`（联表）。

## Model

```go
type SysPost struct {
    PostId int `gorm:"primaryKey;autoIncrement" json:"postId"`
    // ... 业务字段
    models.ControlBy   // CreateBy / UpdateBy
    models.ModelTime   // CreatedAt / UpdatedAt / DeletedAt
}
func (SysPost) TableName() string { return "sys_post" }
```

`TableName()` 必须显式声明（GORM 配置了 `SingularTable`，不会自动推导复数）。

## 公共契约面

第三方应用（`app/` 下的业务模块）可以稳定依赖哪些包、路由与迁移怎么注册、
哪些约束是硬的，见 `docs/contract.md`。

两条与主仓贡献者直接相关的：

- **`common/`、`core/` 不得 import `app/`** —— `make checksilent` 在 CI 里守着，违反即红。
- **从 core 契约包声明出来的类型必须写成别名**（`type X = pkg.Y`，不是 `type X pkg.Y`）
  —— `contract-shim-alias` 检查守着。defined type 会丢掉整个方法集，
  而且**不一定在本仓编译失败**，理由见 `docs/contract.md` 末节。
- **注册类 API（`AppRouters` / `sdk.Runtime.SetAppRouters` / `migration.ForApp`）
  必须在 `runStartupHooks()` 之前调用完** —— `init()` 是最省事的位置，
  但约束的是**顺序**，不是写在哪个函数里；晚到的注册会被丢弃并只记一条 ERROR。

## 路由注册

通过 `init()` 自注册，不在中心文件手工添加：

```go
func init() { routerCheckRole = append(routerCheckRole, registerSysPostRouter) }

func registerSysPostRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
    api := apis.SysPost{}
    r := v1.Group("/post").
        Use(authMiddleware.MiddlewareFunc()).
        Use(middleware.AuthCheckRole()).      // Casbin 鉴权
        Use(actions.PermissionAction())       // 注入数据权限
    { r.GET("", api.GetPage); r.POST("", api.Insert); /* ... */ }
}
```

新增路由文件后，需确认 `cmd/api/` 中已用 `_` 导入该包。

## 命名

| 对象 | 规则 | 示例 |
|---|---|---|
| 数据表 | `sys_` 前缀 + 下划线 | `sys_post` |
| API 路径 | `/api/v1/` + kebab-case | `/api/v1/sys-user` |
| DTO | `{Model}{Action}Req` | `SysPostPageReq` |
| 权限标识 | `模块:资源:操作` | `admin:sysPost:add` |

权限标识需与前端 `v-permisaction` 一致，并写入 `sys_menu` 种子数据——完整可运行的
参照见 `cmd/migrate/migration/version/1786700001000_demo_menu.go`（sys_api /
sys_menu / sys_menu_api_rule / casbin_rule 四张表如何配齐，用的是幂等 upsert，
可以直接照抄结构）。

## Swagger

Handler 必须带完整注解，`go generate` 会据此生成文档：

```go
// @Summary 岗位列表
// @Tags 岗位
// @Success 200 {object} response.Response
// @Router /api/v1/post [get]
// @Security Bearer
```

## 本地运行

**配置 `driver: sqlite3` 时必须带构建标签**，否则启动即 panic：

```bash
go run -tags sqlite3 . migrate -c config/settings.sqlite.yml
go run -tags sqlite3 . server  -c config/settings.sqlite.yml
```

原因：`common/database/open.go` 带 `//go:build !sqlite3`，不加标签时编进的是
不含 sqlite3 的版本，`opens["sqlite3"]` 为 nil，调用时在 nil 函数上崩溃。
报错信息不会提到构建标签，容易误判成环境损坏。MySQL / PostgreSQL 无此问题。

对应 `Makefile` 的 `build-sqlite` 目标。

## 数据库迁移

文件名前 13 位为毫秒时间戳版本号，不合规的名字会在启动时 panic 并报出该文件名。
**已执行过的迁移文件不可修改** ——
`sys_migration` 表按版本号去重，改动不会重跑，只能新增一个迁移来修正。

放哪个目录取决于身份：

| 目录 | 用途 | 是否入库 |
|---|---|---|
| `version/` | 框架自带迁移，随仓库分发给所有使用者 | 是 |
| `version-local/` | 使用者自己项目的迁移 | 否（已在 `.gitignore`） |

**向本仓库提交迁移必须放 `version/`** —— 放进 `version-local/` 会被忽略掉，
`git status` 看不到，PR 里也不会出现。两个目录的包名分别是 `version` 与
`version_local`（后者与目录名不一致，因为标识符不能含连字符）。

### 写种子数据用哪个 models 包

`1786700003000` 之后新增的迁移，**种子数据要用 `app/` 下的运行时模型**
（如 `app/admin/models.SysApi`、`SysMenu`），**不要用 `cmd/migrate/migration/models`**。

后者的 `ModelTime` 声明的是可空的 `gorm.DeletedAt`，这对它之前的迁移是对的（那正是
当时列的形状），转换之后就不再成立，两个方向都会出问题：

- **写**：往 NOT NULL 列里塞 NULL，第一条 insert 就 `NOT NULL constraint failed`
- **读**：GORM 拼 `WHERE deleted_at IS NULL`，而活跃行存的是 `0`，静默查不到——
  照抄 `demo_menu.go` 的授权段落会因此跳过授权，菜单建好、权限没授、迁移仍记为成功

干净库跑不出这个问题，今天所有用该包的迁移都排在转换之前。完整推导见
`schema_coverage_test.go` 里 `TestPostConversionMigrationsAvoidFrozenSeedModels`
的注释，那个测试也守着这条边界。

## 静默失败校验

`make checksilent` 检查七类**不报错、不记日志、行为悄悄变得不对**的问题，
CI 会跑，命中 ERROR 即失败：

| 检查 | 级别 | 静默后果 |
|---|---|---|
| `modeltime-mix` | ERROR | 两个 `ModelTime` 混用，整张表查不到数据 |
| `menu-sort-overflow` | ERROR | 菜单 `sort` 超 127，MySQL tinyint 拒绝写入，迁移中断 |
| `config-value-truncation` | ERROR | `sys_config.config_value` 超 255 字符被静默截断 |
| `menu-id-collision` | ERROR | 两个模块硬编码同一菜单 ID，互相覆盖 |
| `contract-import-boundary` | ERROR | 契约包 import `app/`，应用无法独立编译 |
| `contract-shim-alias` | ERROR | 契约薄壳写成 defined type 而非别名，方法集丢失，本仓可能照常编译、第三方应用编译不过 |
| `menu-name-mismatch` | WARN | 菜单名与前端组件 `name` 不一致，keep-alive 缓存静默失效 |

最后一条要跨仓库比对，只能做正则启发式，因此是 WARN，**不影响退出码**，
且默认跳过；要跑它得指定前端目录：

```bash
make checksilent UI_DIR=../go-admin-ui/src
```

升级门槛：连续 2 个发版周期零误报后转为 ERROR。

## 提交规范

格式 `type+emoji: 描述`：

`feat✨` `fix🐛` `style💄` `docs📝` `perf👌` `test✅` `refactor🎨` `chore🔧`

一个提交只做一件事。改动跨越多个语义时拆分提交，不要混在一起。

## 红线

- 不使用全局 DB 变量，一律用 `e.Orm`（来自请求上下文，多租户依赖它）
- 不在 Service 中引用 `gin.Context`
- 生产部署前确认 `mode: prod` 且已修改 `jwt.secret`（dev 模式下 token 几乎不过期）
- 不提交 `config/settings.yml` 中的真实凭据
