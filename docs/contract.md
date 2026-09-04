# 公共契约面

> 本文写给**第三方应用作者**：你写一个装进 go-admin 的业务模块，可以依赖什么、
> 怎么接进来、哪些约定不遵守会**不报错地出错**。
>
> 主仓贡献者的编码约定见根目录 `AGENTS.md`，设计取舍见 `docs/architecture.md`。

---

## 契约面在 core，不在 go-admin

这份文档以前列的是 go-admin 自己的四个包（`common/actions` 等），依据写的是
「把 `app/demo` 的 import 去重之后恰好就是这四个」。

**那个依据是错的，而且错的方向是把人引向依赖宿主。**

go-admin 的使用方式是 clone / fork：每个使用者拿到的是一整份代码，然后**改它**。
应用如果依赖 `go-admin/common/actions`，它依赖的是一个**每个使用者都不一样、
而且随时在变**的东西——你没有办法测试自己的应用在别人改过的 fork 上能不能编译。

还有一条更硬的：`go-admin` 这个 module path 没有点号，
按 Go 的规则**不是合法的可解析模块路径**：

```
$ go get go-admin/common/models
go: malformed module path "go-admin/common/models": missing dot in first path element
```

想 import 它就必须写 `replace`，而**非主模块的 `replace` 会被忽略**——
你在自己应用里写的 replace 对使用者不生效。所以「应用 require go-admin」
这条路不是不优雅，是走不通。

契约面因此落在 **go-admin-core**：那是唯一一个大家都一样、有版本号、
不会被使用者随手改的东西。

---

## 承诺稳定的包

全部在 `github.com/go-admin-team/go-admin-core/v2` 下：

| 包 | 用途 |
|---|---|
| `sdk/contract/models` | `Model` / `ControlBy` / `ModelTime` / `ActiveRecord` / `BaseUser` / `Migration`、`sys_menu.menu_type` 的三个枚举值 |
| `sdk/contract/dto` | `Pagination` / `MakeCondition` / `Paginate` / `OrderDest` / `ObjectById`、`Index` 与 `Control` 接口 |
| `sdk/contract/actions` | 数据权限设施：`DataPermission` / `Permission` / `PermissionAction` / `GetPermissionFromContext`、五个 `DataScope*` 常量与 `IsValidDataScope` |
| `sdk/contract/migration` | `Registry` / `AppRegistrar` / `ForApp` / `SetVersion` / `GetFilename` |
| `sdk/contract/seed` | `MenuSpec` / `ApiSpec` / `Seeder` / `SeedMenus`——往侧边栏和接口表里登记自己 |
| `sdk/pkg` | `GetOrm(c)`：从请求上下文取本租户的数据库连接 |
| `sdk/api`、`sdk/service` | 可选的 Api / Service 基类 |
| `response` | `OK` / `Error` / `PageOK`：响应格式 |
| `jwtauth/user` | 从 token 取当前用户身份 |
| `sdk/runtime` | 中间件 key 常量与 `GetHandlerFunc`：复用宿主已注册的鉴权链 |

`sdk/contract/` 这个前缀的含义就是「**承诺对应用稳定**的那一面」。core 里
`sdk/` 下的其他包是框架基础设施，语义不同——上表逐个列了名字，
**不要因为「都在 core 里」就认为是契约面**。

"稳定"的含义：**在 core 的 `v2.x` 内不做破坏性变更**。新增导出符号不算破坏；
改签名、改语义、删除导出符号算，会走 major 版本并在 release note 里单列。

准确的语义以 core 那份文档为准：
[go-admin-core `docs/contract.md`](https://github.com/go-admin-team/go-admin-core/blob/main/docs/contract.md)。
本文写的是宿主这一侧——它管不着的那些。

### go-admin 自己的包

`go-admin/common/models`、`common/dto`、`common/actions` 里的契约类型现在是
**指向 core 的类型别名**（`type X = corepkg.X`），主仓和所有 fork 的存量代码
一行不用改。别名在编译期就是同一个类型，不是"兼容层"。

但**新写的应用不要 import 它们**——那样就又依赖上宿主了。

---

## 契约面是三层，不是一层

划分依据不是"应用会 import 哪些包"，而是**"哪一条不遵守会静默出错"**：

| 层 | 内容 | 判据 |
|---|---|---|
| **一 · 必须遵守** | 路由注册、从 context 取库、响应 shape、`ControlBy`/`ModelTime`、鉴权、数据权限、事务范式 | 不遵守 → **不报错，行为悄悄不对** |
| **二 · 可选便利** | `api.Api`、`service.Service`、CRUD Action、`MakeCondition` | 用不用都对 |
| **三 · 今天空白** | 应用间调用、领域事件、缓存租户隔离 | **没有。别自己发明** |

**框架不强制任何一层抽象。** 一个不用任何便利层的 handler 完全合法：

```go
func handler(c *gin.Context) {
    db, err := pkg.GetOrm(c)
    if err != nil {
        response.Error(c, 500, err, "")
        return
    }
    var list []MyModel
    if err := db.Find(&list).Error; err != nil {
        response.Error(c, 500, err, "")
        return
    }
    response.OK(c, list, "")
}
```

第一层则是不管你用不用便利层都要遵守的，逐条写在下面，每条都附**不遵守会怎样**。

---

## 第一层：不遵守就静默出错

### 1. 路由注册

见下方「注册路由」一节。

**不遵守会怎样**：注册表在 `RunAppRouters()` 之后就封闭了，晚到的注册被丢弃，
只记一条 ERROR 日志。包级 `AppRouters` 连这个都没有——它就是一个普通 slice，
什么时候 append 都"成功"，启动钩子之后 append 的那些永远不会执行，**且不出声**。

### 2. 数据库连接从 context 取，不用全局变量

```go
db, err := pkg.GetOrm(c)          // 唯一正确的取法
```

`common/middleware/db.go` 在每个请求上按 `c.Request.Host` 挑出本租户的连接
放进 context：

```go
c.Set("db", sdk.Runtime.GetDbByTenant(c.Request.Host).WithContext(c))
```

**不遵守会怎样**：连接是**按租户注册**的（`SetDbByTenant(host, db)`），
`GetOrm(c)` 按 `c.Request.Host` 挑。你要是在启动时把某个连接存进包级变量再一直用，
多租户部署下所有租户的读写就都落到那一个库上——不报错、不告警，数据串了才发现。

这个坑在本仓库真踩过：`common/global.Driver` 取的是启动循环
**迭代到的第一个**库的驱动（`common/database/initialize.go`），
而 Go 的 map 迭代顺序是随机的——两个库用不同驱动时，那个值每次启动都可能不一样。
所以「一个进程一个库」这个假设不要写进任何一行代码。

### 3. 响应 shape

一律用 `response.OK` / `response.Error` / `response.PageOK`，不要自己
`c.JSON`。它们发出去的形状是：

```jsonc
// 成功
{"requestId": "...", "code": 200, "data": {...}}
// 分页：data 里再套一层
{"requestId": "...", "code": 200, "data": {"count": 42, "pageIndex": 1, "pageSize": 10, "list": [...]}}
// 失败
{"requestId": "...", "code": 500, "msg": "...", "status": "error"}
```

**HTTP 状态码永远是 200**，业务码在 body 的 `code` 里——这是既定行为，
`response.Error` 走的是 `c.AbortWithStatusJSON(http.StatusOK, res)`。

**不遵守会怎样**：前端 `src/utils/request.ts` 的响应拦截器只读 body 的 `code`，
`code !== 200` 就弹一条 `msg` 内容的 error toast 并 reject。你自己
`c.JSON(200, myThing)` 的话 `code` 是 `undefined`，界面上弹出来的是**一条空的
错误提示**，数据到不了页面。列表更安静：`useTable.ts` 读的是
`page?.list ?? []` 和 `page?.count ?? 0`，形状对不上就是**一张空表，零报错**。

### 4. `ControlBy` 与 `ModelTime`

每张业务表的 model 都嵌这三个：

```go
type Order struct {
    models.Model        // Id
    // ... 你的字段 ...
    models.ControlBy    // CreateBy / UpdateBy
    models.ModelTime    // CreatedAt / UpdatedAt / DeletedAt
}

func (Order) TableName() string { return "app_order" }   // 必须显式声明
```

`ControlBy` 提供 `create_by` 列，**数据权限的每一条 SQL 都 join 在它上面**。
`ModelTime` 的 `DeletedAt` 是 `soft_delete.DeletedAt`（毫秒时间戳，活行为 0，
永不为 NULL），不是 `gorm.DeletedAt`。

**不遵守会怎样**：

- 嵌了 `ControlBy` 但写入时忘了 `SetCreateBy(user.GetUserId(c))`，
  `create_by` 就是 0。除「全部数据权限」外的每一档都**查不到任何数据**，
  而且不报错——看起来像"这个用户还没建过数据"。
- 用错 `ModelTime` 版本（可空的 `gorm.DeletedAt`）：gorm 按
  `deleted_at IS NULL` 过滤，而活行里存的是 0，于是**整张表一行都查不出来**。
  主仓的 `sys_columns` / `sys_tables` 真在这个状态下待过——代码生成器
  一张表都列不出来，没有任何报错。`make checksilent` 的 `modeltime-mix`
  就是为这条加的。
- `TableName()` 忘了写：GORM 配了 `SingularTable`，不会推导复数，表名会是
  你没预料的那个。

### 5. 鉴权：用宿主已注册的中间件，不要自己造

```go
jwtCheck, ok := sdk.Runtime.GetHandlerFunc(runtime.JwtTokenCheck)
if !ok {
    log.Fatal("JwtTokenCheck is not registered; is the host started via cmd/api?")
}
roleCheck, _ := sdk.Runtime.GetHandlerFunc(runtime.RoleCheck)
permCheck, _ := sdk.Runtime.GetHandlerFunc(runtime.PermissionCheck)

g := v1.Group("/order").Use(jwtCheck).Use(roleCheck).Use(permCheck)
```

三个 key 的常量在 `sdk/runtime`，宿主启动时把三个中间件注册进去。

**不遵守会怎样**：`GetHandlerFunc` 在"没注册"和"注册成了别的类型"两种情况下
都返回 `ok=false` 而不是 panic——**因为路由注册跑在 core 的 panic 护栏里面，
裸类型断言 panic 之后日志报的是"这个模块一条路由都没注册上"，跟真实原因对不上**。
所以 `ok` 必须自己判，判出来要**大声失败**：一个跳过鉴权继续注册的路由，
就是一条静默的匿名可访问接口。

**宿主必须注册绑定过的闭包。** 三个 key 存的都得是 `gin.HandlerFunc`——
比如 `authMiddleware.MiddlewareFunc()`，**不是** `(*jwt.GinJWTMiddleware).MiddlewareFunc`。
后者是方法表达式，没有接收者绑在上面，取回来断言不成 `gin.HandlerFunc`，
怎么断言都做不成一个能用的 handler。

> **当前状态**：`common/middleware/init.go` 里 `RoleCheck` 与 `PermissionCheck`
> 注册的是 `AuthCheckRole()` 和 `actions.PermissionAction()`，都是绑定过的闭包，
> 取回来就能用；**`JwtTokenCheck` 注册的还是那个方法表达式**，所以今天对它
> `GetHandlerFunc` 拿到的是 `ok=false`。上面那段 `log.Fatal` 会在启动时打出来——
> 这是有意的，宁可起不来也不要一条没鉴权的路由。主仓这一处的修复见 F10，
> 修完之后本段可以删掉。

还有一条**不影响行为但影响理解**的：主仓今天四个模块各自调一次 `AuthInit()`
（`app/admin`、`app/jobs`、`app/other`、`app/demo`），也就是有四个 JWT 实例。
这不产生行为差异——配置同源（`config.JwtConfig`），JWT 校验是无状态的，
不看实例身份。但它意味着 `GetHandlerFunc(runtime.JwtTokenCheck)` 取回来的是
**最后注册进去的那一个**。要让应用拿到一个有意义的共享实例，宿主应当在注册路由
之前构造一次，而不是每个模块构造一次。

**测的时候别用 `admin` 账号。** `AuthCheckRole` 里 `rolekey == "admin"` 直接
`c.Next()`，**完全跳过 Casbin**。拿 admin 压任何鉴权路径都测不到东西。

### 6. 数据权限

两件事都要做：

```go
// 路由上挂中间件（上一节的 permCheck 就是它）
g := v1.Group("/order").Use(permCheck)

// 查询里组合 scope
p := actions.GetPermissionFromContext(c)
db.Scopes(actions.Permission(Order{}.TableName(), p)).Find(&list)
```

`sys_role.data_scope` 有五档，`Permission()` 按它拼 WHERE 条件：

| 值 | 常量 | 含义 | 条件 |
|---|---|---|---|
| `1` | `DataScopeAll` | 全部数据权限 | 不加条件 |
| `2` | `DataScopeCustom` | 自定义数据权限 | `create_by` 属于 `sys_role_dept` 关联到的部门 |
| `3` | `DataScopeDept` | 本部门 | `create_by` 属于本部门 |
| `4` | `DataScopeDeptTree` | 本部门及以下 | `create_by` 属于 `dept_path` 匹配的子树 |
| `5` | `DataScopeSelf` | 仅本人 | `create_by = 当前用户` |

自己往 `sys_role.data_scope` 写值的话先过一遍 `IsValidDataScope`——
写进去的非法值不会在写入时报错，只会在**每一次查询**里静默地什么都查不到。

**不遵守会怎样**，两种漏法的方向相反，值得分清：

- **查询里忘了组合 `Permission()`** —— 就是**全量可见**，每个角色都看得到所有人
  的数据，不报错、不记日志。**这是本框架里最贵的一类静默失败**，所以那一行
  `db.Scopes(...)` 不是"最佳实践"，是契约。
- **组合了 `Permission()` 但路由上漏挂中间件** —— 上下文里没有 `PermissionKey`，
  拿到的是零值，`DataScope` 是空串，落进下面那个 fail-closed 的 default，
  结果是**一行都查不到**。方向反了，至少还看得见。

五档之外的值（空串、拼错的、还没迁移的老数据）落到 `default` 分支，
那里是 **fail closed**：加一条 `1 = 0`，什么都不返回。注意 `1`（全部数据权限）
是**显式列出的一个 case**，不是"落到 default"——两者曾经是同一条路，
于是"没配置"和"配置成看全部"产出的 SQL 一个字都不差。

`3` / `4` 两档在 `DeptId <= 0` 时同样 fail closed。原因是
`sys_dept.dept_path` 一律以 `/0/` 开头，`dept_id=0` 会把 LIKE 模式变成
`'%/0/%'`，**命中全表**——本来想表达"没有部门"，实际表达的是"全部部门"。

数据权限还有一个**全局开关** `application.enabledp`，默认是 `false`。
关掉时 `Permission()` 原样返回查询、`PermissionAction()` 直接放行——
**你的应用在默认配置下测不出数据权限的任何行为**，要验证得先把它打开。

**不要自己重写这段 SQL。** 那 20 行里埋着 8 项内部知识：JWT claims 的私有键名
（`datascope` / `deptid`）、`sys_user`↔`sys_role` 的 join、`sys_role_dept`
关联表、`sys_dept.dept_path` 的 `/0/1/2/` 编码、`create_by` 的归属约定、
`enabledp` 开关、老 token 的回落逻辑。**而且写错的方向是越权。**
仓库里有过一份第二实现，`dept_path` 的匹配写成 `"%"+id+"%"` 少了两个斜杠，
`dept_id=1` 会匹配上 `/11/`、`/21/`、`/100/`——写它的人比第三方更懂这套约定，
仍然写错了。那份实现已经删掉了。

### 7. 事务范式

**业务层的事务一律用 `Transaction()` 闭包形式**：

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&order).Error; err != nil {
        return err   // rolled back
    }
    return tx.Model(&stock).Where("qty >= ?", n).
        UpdateColumn("qty", gorm.Expr("qty - ?", n)).Error
})
```

GORM 自己处理提交、回滚，以及 **panic 时的回滚**。

**不要照抄 `app/admin/service/sys_role.go`。** 那里有 5 处手写的
`Begin` / `defer` 写法，三个缺陷都是静默的：

```go
tx := e.Orm
if config.DatabaseConfig.Driver != "sqlite3" {   // 缺陷 2
    tx = e.Orm.Begin()
    defer func() {
        if err != nil { tx.Rollback() } else { tx.Commit() }   // 缺陷 1
    }()
}
```

1. **panic 时提交半截事务**——defer 只看 `err`，panic 时 `err` 仍是 nil，走的是
   `Commit()`
2. **sqlite 下根本不开事务**——那一整个特判让 `tx` 就是 `e.Orm` 本身，
   写一半失败留一半
3. **读 `config.DatabaseConfig.Driver`**——那是全局单库配置，多租户下不是
   当前租户的驱动

缺陷 1 不止那一处：`app/admin/service/sys_dept.go`、`sys_menu.go`、
`app/other/models/tools/sys_tables.go` 用的是同一个 `defer` 写法
（没有 sqlite 特判，所以只有缺陷 1）。**整个 `Begin`/`defer` 家族都别照抄。**

同一个仓库里就有正确的参照：`cmd/migrate/migration/version/` 下 7 个迁移里
5 个用的是闭包形式（另外两个是纯 DDL 标记，DDL 在 MySQL 下本来就不进事务），
且这条路在 sqlite 下实测跑得通（`make build-sqlite`）。
主仓那些写法本批次不改，单独跟。

**并发保护用条件更新 + `RowsAffected`**，不要"先查后改"：

```go
res := tx.Model(&Order{}).Where("id = ? AND status = ?", id, StatusPending).
    Update("status", StatusPaid)
if res.Error != nil { return res.Error }
if res.RowsAffected == 0 { return ErrAlreadyPaid }   // 别人先改了
```

---

## 第二层：可选便利

用不用都对，**不用不会出任何问题**：

| 东西 | 在哪 | 是什么 |
|---|---|---|
| `api.Api` | core `sdk/api` | 一条链式糖：`MakeContext` / `Bind` / `MakeOrm` / `OK` / `PageOK` / `Error` |
| `service.Service` | core `sdk/service` | 一个装 `Orm` / `Log` / `Cache` / `Error` 的结构体加一个 `AddError` |
| `MakeCondition` / `search` tag | core `sdk/contract/dto` | 把 DTO 上的 `search:"type:exact;column:name;table:xx"` 翻成 WHERE |
| 通用 CRUD Action | go-admin `common/actions` | `IndexAction` 等五个。**留在 go-admin，没有下沉** |

最后一行是有意的：CRUD Action 是最需要演进的一类东西（分页参数、批量操作、
软删语义、字段级权限），而 core 的每一个导出都是永久承诺——放进去容易，
拿出来不可能。想用就把那 294 行抄走，抄走的那份还能按你自己的需要改。
主仓唯一的真实业务模块 `app/admin` **一个 CRUD Action 都没用**，全是手写 Service。

`MakeCondition` 返回的是 `func(db *gorm.DB) *gorm.DB` 闭包，方言从闭包里那个
`db.Dialector.Name()` 读，**必然是本租户那个库的驱动**，不需要你设置任何东西。

---

## 第三层：今天没有的

**明说没有，别自己发明**：

| 能力 | 现状 |
|---|---|
| 应用间调用 | 零定义。A 应用要调 B 应用只能直接 import 对方的包，循环依赖就回来了 |
| 领域事件 / EventBus | 无 |
| 缓存的租户隔离 | `service.Service` 有 `Cache` 字段，**是否按租户隔离未验证**。当作没隔离来写 |
| 异步任务 | 有队列，但热更新后消费者会丢（issue #892） |

这几条留给后续批次，按真实需求补——现在凭空设计只会设计错。
如果你的应用卡在这里，在 issue 里说一声，那正是我们要的输入。

---

## 装一个应用要接两处线

后端**两处**，漏掉第二处是**静默失败**：

```go
// 1. 路由：cmd/api/<name>.go
import _ "github.com/acme/go-admin-app-order/router"

// 2. 迁移：cmd/migrate/server.go 的 import 块里
import _ "github.com/acme/go-admin-app-order/migration"
```

两个都是空导入，作用只是让那个包的 `init()` 跑起来。

**漏了第二处会怎样**：不报错。`migrate` 命令照常跑完、照常打印成功，
你的建表和种子数据**就是不执行**。等到第一个请求打过来才会看到
"表不存在"，而那时排查方向已经跑偏了。

`migrate --dry-run` 是确认接线成功的最快方式——它只读，可以直接对生产库跑：

```bash
go-admin migrate --dry-run -c config/settings.yml    # 你的迁移应该出现在列表里
```

带界面的应用还有第三处，在前端仓库，见下一节。

---

## 前端：菜单 `component` 必须以 `apps/` 开头

前端那一处接线是 `go-admin-ui` 的 `apps.config.mjs`——加一条
`{ code: 'order', source: '...' }`，`source` 指到你的页面目录
（兄弟目录的相对路径，或 `./node_modules/@scope/app-order/views/order`）。
`scripts/sync-apps.mjs` 会在 `pnpm dev` 与 `pnpm build` 之前把它复制进
`src/apps/<code>/`，不需要手工跑。

`src/stores/permission.ts` 的 `appPath()` **只认路径第一段是 `apps`**，
其余一律当成主仓内置视图去 `src/views/` 下找。

所以你的菜单种子里 `Component` 必须写成：

```
apps/<code>/<该应用内的相对路径>/index
```

比如 `code` 是 `order` 的应用写 `apps/order/index`（开头带不带 `/` 都行，
只看第一段）。**不能**写成 `/order/index`。

**写错会怎样**：第一段是 `order` 而不是 `apps`，前端会去找一个不存在的
`src/views/order/index.vue`，页面摔到 `AppNotInstalled` 占位组件。
但控制台打印的是 `no component at src/views/order/index.vue`——
**跟真实原因（漏了 `apps/` 前缀）对不上**，排查时很容易被这条日志带偏。

对应的前端约定写在 go-admin-ui 的 `AGENTS.md`。另外一条：`source` 目录的内容
**原样**搬进 `src/apps/<code>/`，不会在 `code` 之外再自动插一层——想要
`apps/order/index` 这种最短形式，`source` 就要直接指到该应用**这一个页面模块**
的目录，而不是应用仓库的 `views` 根目录。

---

## 注册路由

写一个 `func()` 签名的 `InitRouter`（照抄 `app/demo/router/router.go`），
然后二选一接进来：

```go
// 方式一（历史写法，仍然有效）：在主仓 cmd/api/<name>.go 里
AppRouters = append(AppRouters, router.InitRouter)

// 方式二（推荐）：不需要 import go-admin/cmd/api
sdk.Runtime.SetAppRouters(router.InitRouter)
```

**第三方应用只能走方式二**——方式一要求 `import "go-admin/cmd/api"`，
那就又依赖上宿主了。

走方式二还多拿到两样东西，都在 core 那边实现：

- **panic 护栏**——你的 `InitRouter` panic 了，其余模块照常注册、进程不退出，
  日志里会写明是哪一行注册的
- **失败分级**——`sdk.Runtime.SetAppRoutersWith(f, runtime.WithFatal())`
  声明「我起不来就别启动」

方式一（包级 `AppRouters`）没有护栏，panic 直接掀桌。

**执行顺序**：先跑完包级 `AppRouters`，再由 `sdk.Runtime.RunAppRouters()`
跑 core 自己的注册表，各自内部保持注册顺序。别依赖跨来源的相对顺序。

`InitRouter()` 内部：自己拿 `sdk.Runtime.GetEngine()`，按需建
`gin.RouterGroup`，通过 `init()` 自注册到你自己包内的
`routerCheckRole` / `routerNoCheckRole` 列表，不在任何中心文件手工列举。

---

## 注册数据库迁移

框架自身的迁移用 `SetVersion`；应用的迁移走 `ForApp`：

```go
func init() {
    _, fileName, _, _ := runtime.Caller(0)
    migration.ForApp("crm").SetVersion(migration.GetFilename(fileName), initCrmTables)
}

func initCrmTables(db *gorm.DB, version, appCode string) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // ... schema / data changes ...
        return tx.Create(&models.Migration{Version: version, AppCode: appCode}).Error
    })
}
```

注册面（`ForApp` / `SetVersion` / `GetFilename`）在 core 的
`sdk/contract/migration`，是一个**进程级的包级注册表**——`ForApp` 直接当包级函数
调，不需要从宿主手里接过什么句柄。**执行面**——读 `sys_migration`、排序、跑事务、
`migrate` 与 `migrate status` 两个命令——留在宿主，它通过 `Snapshot()` 读那张表。

仓库内的模块继续经 `go-admin/cmd/migrate/migration` 走，那个包现在是薄壳，
导入路径不变；外置应用直接 import core 的那个包，**两边写法一模一样**。

四条必须知道的规则：

1. **完成记录由迁移函数自己写**，而且要写在自己的事务里。框架的调度循环只做
   "这个 version 在 `sys_migration` 里有没有" 的判断，从不代你插入 —— 这样
   "数据改完了"和"标记成已完成"才是同一个事务，不会出现改了一半却被记成成功。
2. **`AppCode` 必须写进去**。签名多带一个 `appCode` 参数就是为此 —— 忘了写，
   schema 上那一列等于白加，你的迁移会被记成框架的。
3. **落库的 `version` 是加了前缀的**。`ForApp("crm")` 注册 `1786800001000`，
   实际写进 `sys_migration.version` 的是 `crm-1786800001000`，函数收到的
   `version` 参数已经是这个带前缀的值，照抄进 `models.Migration{Version: version}`
   即可。前缀的意义是：两个来源不同的应用哪怕碰巧生成同一个毫秒时间戳，也不会撞主键、
   不会有一方被误判为"已应用"。
4. **应用 code 一律小写**，`ForApp` 会自己 `strings.ToLower` 一遍。`core` 是保留字
   （`migrate status` 用它表示框架自身，`--app core` 选中框架），`ForApp("core")`
   会 panic。

顺序保证：**同一应用内按版本号严格有序**。跨应用顺序不做承诺 —— 由于前缀的存在，
今天的实际顺序是"先跑完全部框架迁移，再按 appCode 字母序逐个应用跑完"，
但这是实现细节，不要依赖它。跨应用依赖（应用 A 的迁移要求应用 B 先跑完）
需要依赖拓扑排序，属于后续阶段。

看当前状态、看这次会跑什么，不用猜：

```bash
go-admin migrate status -c config/settings.yml          # 按应用分组列出已应用 / 待应用
go-admin migrate --dry-run -c config/settings.yml       # 列出会执行什么、什么顺序，不写库
go-admin migrate --app crm -c config/settings.yml       # 只跑 crm 的迁移
```

`status` 与 `--dry-run` 是纯只读的，不建表、不改表结构，可以直接对生产库执行。

---

## 菜单与接口种子

一个带界面的应用要在侧边栏里出现，需要往四类数据里写东西：`sys_api`、
`sys_menu`、`sys_menu_api_rule`（菜单与接口的关联）、以及角色授权与 Casbin
策略（`sys_role_menu` / `casbin_rule`）。

**你不需要知道这些表长什么样。** `sdk/contract/seed` 让你只描述"我要什么"，
由宿主决定"怎么写进它自己的表"：

```go
// 在你自己的迁移里，用它自己的那个事务
err := seed.SeedMenus(tx, "order", []seed.MenuSpec{
    {Code: "root", Kind: models.Directory, Title: "订单"},
    {Code: "list", Parent: "root", Kind: models.Menu, Title: "订单列表",
        Path: "/order", Component: "apps/order/index", ApiCodes: []string{"list"}},
}, []seed.ApiSpec{
    {Code: "list", Title: "订单列表", Path: "/api/v1/order", Method: "GET"},
})
```

`Kind` 用的就是 `sdk/contract/models` 里 `sys_menu.menu_type` 的那三个值
（`Directory` / `Menu` / `Button`），不是另一套同值的常量。

`Component` 的写法见上面「前端」一节——**这里是最容易写错的一个字段**。

core 里**没有** `SysMenu`、没有 `SysApi`、没有任何表名。这是刻意划的边界：
这个框架的宿主里本来就已经有两份 `SysMenu`（一份冻结在迁移期、一份运行期），
两者在软删语义上不一致，害过人，为此专门建了一个仓库内的工具来守。
往 core 里再放第三份表结构，就等于在**唯一没有工具守着**的地方重造同一类 bug。

---

## 应用配置节

不要改宿主的源码去加配置。`sdk/config.RegisterExtend` 让你认领
`extend:` 下自己那一节：

```go
type orderConfig struct {
    PaymentEndpoint string
    Timeout         int
}

// 在 init() 里调，与 SetAppRouters / ForApp 同一约定
var getOrderConfig = config.RegisterExtend[orderConfig]("order")

func handler(c *gin.Context) {
    cfg := getOrderConfig()
    _ = cfg.PaymentEndpoint
}
```

```yaml
extend:
  order:
    PaymentEndpoint: https://payment.internal
    Timeout: 30
```

每个 key 各自解码，互不覆盖。**同一个 key 注册两次会立刻 panic**——
注册期没有"封闭时刻"可以用来拒绝迟到的注册，所以重复只能在注册的那一刻
大声报出来，而不是让第二个人静默顶掉第一个人的配置节。

配置文件是被监听的，改动会触发重载。`RegisterExtend` 每次重载解码进一个全新的
`T` 再原子换指针，所以访问器拿到的永远是一个自洽的快照，请求路径上读它不需要加锁。
唯一要注意的：**不要跨两次调用拼一个视图**——从同一个返回值上读两个字段是一致的，
调两次访问器各读一个字段，中间夹一次重载就不是了。

---

## 硬约束：注册要赶在启动钩子之前

三个注册入口——`AppRouters`、`sdk.Runtime.SetAppRouters`、`migration.ForApp`——
都必须在 `cmd/api/server.go` 的 `runStartupHooks()` 执行之前调用完。

`init()` 是最省事的位置：Go 规范保证包级变量初始化与 `init()` 在 `main()` 之前
**单 goroutine 顺序执行**，注册期天然没有并发写。但它不是唯一合法位置——
在 `run()` 之类早于启动钩子的地方注册同样成立。**这条规则约束的是顺序，
不是你写在哪个函数里。**

想在代码里判断注册窗口是否还开着：

```go
if sdk.Runtime.AppRoutersSealed() { /* RunAppRouters 已经跑过了 */ }
```

主仓这边补三条 core 那份文档管不着的：

1. **`RunAppRouters()` 跑过之后，core 的注册表就封闭了**，再调
   `sdk.Runtime.SetAppRouters` 会被丢弃并记一条 ERROR 日志。包级 `AppRouters`
   没有这个机制——它就是一个普通 slice，什么时候 append 都"成功"，
   但 `runStartupHooks()` 之后 append 的那些永远不会被执行，且不出声。
   这是继续推荐方式二的理由之一。
2. **封闭是黏性的，而 `sdk.Runtime` 是包级单例。** 写测试时若会触发启动钩子，
   必须换掉它再还原，否则同一个测试二进制里后面的测试会静默丢注册：

   ```go
   previous := sdk.Runtime
   t.Cleanup(func() { sdk.Runtime = previous })
   sdk.Runtime = runtime.NewConfig()
   ```

   `cmd/api/server_test.go` 里的 `freshRuntime` 就是这个。
3. **迁移的调度循环是主仓的东西**，core 只有注册面。迁移注册的约束仍然是
   "赶在调度循环跑起来之前"，实践上就是 `init()`。

---

## 安全边界：装一个应用等于信任它

**这一层划不出安全边界，本文不假装划得出。**

第三方应用的代码在**宿主进程内**运行，与宿主**同权限**。它持有的是裸的
`*gorm.DB`——`seed.SeedMenus` 用的就是你自己迁移里那个 `tx`，绕开 `Seeder`
直写 `sys_menu`、`sys_api`、甚至 `casbin_rule` 一直都做得到，Go 的类型系统
拦不住，本框架的任何一层也拦不住。

还有一条**不碰 `casbin_rule` 也能走通**的间接路径：把自己的菜单通过
`ApiCodes` 关联到别人的接口，然后等管理员在后台把这个菜单授权给某个角色——
策略是后台自己生成的，记在管理员头上。

所以：

> **装一个应用，等于信任它。** 这和 `import _` 一个 Go 库是同一量级的信任。
> `Seeder` 这类设计的目的是让**守规矩的应用不必知道宿主的表结构**，
> 不是把不守规矩的应用关起来。

给使用者的实际建议只有一条：**按信任 Go 依赖的标准来审应用**——看源码、
钉版本、认作者。不要因为它叫"应用"就以为它跑在沙箱里。

---

## 边界由 CI 守着

`tools/checksilent` 里有两条盯契约面的检查，`make checksilent` 在 CI 里跑，
命中 ERROR 即失败：

| 检查 | 盯的是 |
|---|---|
| `contract-import-boundary` | `common/`、`core/` 不得 import `app/`——否则一个删掉 `app/admin` 的 fork 就编译不了它被告知可以依赖的那一面 |
| `contract-shim-alias` | 从 core 契约包声明出来的类型必须是**别名**（`type X = pkg.Y`），不能是 defined type。判据是右手边，不是一份包名清单，所以谁在哪加的都算 |

第二条守的是一条一个字符的差别。`type X = pkg.Y` 和 `type X pkg.Y`
看着几乎一样，但后者只拿走底层结构、**丢掉整个方法集**，于是嵌了它的 model
不再满足 `ActiveRecord`。麻烦在于这**不一定在本仓编译失败**——本仓只用接口
使唤其中一部分类型，没被使唤到的那些在这里编译得好好的，
**到第三方应用或某个 fork 里才炸**，而那里没人看着。

测试文件同样算——一个删掉 `app/admin` 的 fork 也应该能跑 `go test ./...`。

**这两条工具都只扫仓库树。** 装在 module cache 里的第三方应用，
`checksilent` 一个文件都看不到。所以它保的是**这个仓库和它的 fork**，
不是你的应用——你的应用要自己跑自己的检查。

`checksilent` 还检查另外五类"不出声的失败"，写模块时值得先看一眼
`go run ./tools/checksilent -h`。
