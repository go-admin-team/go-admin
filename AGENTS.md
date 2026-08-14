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

权限标识需与前端 `v-permisaction` 一致，并写入 `sys_menu` 种子数据。

## Swagger

Handler 必须带完整注解，`go generate` 会据此生成文档：

```go
// @Summary 岗位列表
// @Tags 岗位
// @Success 200 {object} response.Response
// @Router /api/v1/post [get]
// @Security Bearer
```

## 数据库迁移

新增迁移放 `cmd/migrate/migration/version-local/`，文件名前 13 位为时间戳。
**已执行过的迁移文件不可修改** —— `sys_migration` 表按版本号去重，改动不会重跑。

## 提交规范

格式 `type+emoji: 描述`：

`feat✨` `fix🐛` `style💄` `docs📝` `perf👌` `test✅` `refactor🎨` `chore🔧`

一个提交只做一件事。改动跨越多个语义时拆分提交，不要混在一起。

## 红线

- 不使用全局 DB 变量，一律用 `e.Orm`（来自请求上下文，多租户依赖它）
- 不在 Service 中引用 `gin.Context`
- 生产部署前确认 `mode: prod` 且已修改 `jwt.secret`（dev 模式下 token 几乎不过期）
- 不提交 `config/settings.yml` 中的真实凭据
