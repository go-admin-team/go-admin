# 公共契约面

> 本文写给**第三方应用作者**：你写一个装进 go-admin 的业务模块，可以依赖什么、
> 怎么注册进来、哪些东西随时可能变。
>
> 主仓贡献者的编码约定见根目录 `AGENTS.md`，设计取舍见 `docs/architecture.md`。

---

## 承诺稳定的包

| 包 | 用途 |
|---|---|
| `common/actions` | 通用 CRUD Action（Index / View / Create / Update / Delete / Permission） |
| `common/dto` | 分页、`search` tag 解析、`Control` / `Index` 接口 |
| `common/models` | `ActiveRecord`、`ControlBy`、`ModelTime`、`Model` |
| `common/middleware` | `AuthCheckRole`、`InitMiddleware` 等 |

**依据不是拍脑袋列的**：`app/demo` 是一个可编译、有测试、CI 会跑的标准 CRUD 模块，
把它的 `go-admin/` 前缀 import 全部去重之后，恰好就是这四个包 —— 它代表
"写一个标准模块所需要的最小依赖面"。你的模块如果需要第五个包，先在 issue 里说一声，
那多半意味着契约面缺了什么。

"稳定"的含义：**在 `2.x` 内不做破坏性变更**。新增导出符号不算破坏；改签名、
改语义、删除导出符号算，会走 major 版本并在 release note 里单列。

### 没有已知例外

这四个包**不 import `app/` 下的任何东西**，2026-08-31 起由 CI 强制
（见下方「边界由 CI 守着」）。在此之前有两处反向依赖，都已根治：

| 原位置 | 反向依赖 | 处理 |
|---|---|---|
| `common/middleware/logger.go` | `app/admin/service/dto` 的两个操作日志状态常量 | 常量下沉到 `common/global`，`dto` 侧保留同名常量作为 deprecated 别名，fork 不受影响 |
| `common/middleware/handler/auth.go` | `app/admin/models` 的 `SysUser` / `SysRole` | 该段断言恒失败、设的是零值且开源版无人读取，属死代码，已删除 |

之所以不把它们记成"已知例外"：这份文档的作用就是告诉你哪些包可以依赖，
如果第一条下面就挂着例外脚注，后来人会照着例外抄，边界从第一天起就是脏的。

---

## 其余包不保证稳定

`common/` 下没有出现在上表里的包（`common/global`、`common/storage`、
`common/database`、`common/file_store`、`common/response`、`common/service`、
`common/apis`、`common/middleware/handler`、根 `common` 包……）以及
`app/admin` 的内部实现，**均不承诺稳定**。

其中 `common/global`、`common/middleware/handler`、根 `common` 包是
`common/middleware` 的编译期依赖 —— 它们会被一起拉进你的依赖图，但这不代表
它们的 API 稳定。**不要因为"都在 `common/` 目录下"就认为是契约面。**

规划中的 001（模块路径改名）会把非契约包移进 `internal/`，由编译器强制这条边界。
届时上表之外的包对外部模块直接不可见 —— 现在就照上表写，那次改动对你零成本。

---

## 注册路由

一个应用模块要注册自己的路由，写一个 `func()` 签名的 `InitRouter`
（照抄 `app/demo/router/router.go`），然后二选一接进来：

```go
// 方式一（历史写法，仍然有效）：在主仓 cmd/api/<name>.go 里
AppRouters = append(AppRouters, router.InitRouter)

// 方式二（推荐）：不需要 import go-admin/cmd/api
sdk.Runtime.SetAppRouters(router.InitRouter)
```

方式二是本次新接上的。差别只有一个但很关键：方式一要求你的模块
`import "go-admin/cmd/api"` —— 那是主程序的命令包，让业务模块依赖它很别扭，
也正是"主仓要为每个模块加一个七行文件"的根源。

**执行顺序**：先跑完包级 `AppRouters`，再由 core 的 `sdk.Runtime.RunAppRouters()`
跑它自己的注册表，各自内部保持注册顺序。别依赖跨来源的相对顺序，各模块的
`RouterGroup` 前缀互不相同，本来就不该有顺序依赖。

走方式二还多拿到两样东西，都在 core 那边实现（见
[core 的 `docs/contract.md`](https://github.com/go-admin-team/go-admin-core/blob/main/docs/contract.md)）：
**panic 护栏**——你的 `InitRouter` panic 了，其余模块照常注册、进程不退出，日志里会写明
是哪一行注册的；**失败分级**——`sdk.Runtime.SetAppRoutersWith(f, runtime.WithFatal())`
声明「我起不来就别启动」。方式一（包级 `AppRouters`）没有护栏，panic 直接掀桌。

`InitRouter()` 内部的约定：自己拿 `sdk.Runtime.GetEngine()`，按需建
`gin.RouterGroup`，通过 `init()` 自注册到你自己包内的
`routerCheckRole` / `routerNoCheckRole` 列表，不在任何中心文件手工列举
（与 `AGENTS.md`「路由注册」一节一致）。

---

## 注册数据库迁移

框架自身的迁移不变：

```go
migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700001000DemoMenu)
```

应用的迁移走 `ForApp`：

```go
func init() {
    _, fileName, _, _ := runtime.Caller(0)
    migration.ForApp("crm").SetVersion(migration.GetFilename(fileName), initCrmTables)
}

func initCrmTables(db *gorm.DB, version, appCode string) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // ... schema / data changes ...
        return tx.Create(&common.Migration{Version: version, AppCode: appCode}).Error
    })
}
```

四条必须知道的规则：

1. **完成记录由迁移函数自己写**，而且要写在自己的事务里。框架的调度循环只做
   "这个 version 在 `sys_migration` 里有没有" 的判断，从不代你插入 —— 这样
   "数据改完了"和"标记成已完成"才是同一个事务，不会出现改了一半却被记成成功。
2. **`AppCode` 必须写进去**。签名多带一个 `appCode` 参数就是为此 —— 忘了写，
   schema 上那一列等于白加，你的迁移会被记成框架的。
3. **落库的 `version` 是加了前缀的**。`ForApp("crm")` 注册 `1786800001000`，
   实际写进 `sys_migration.version` 的是 `crm-1786800001000`，函数收到的
   `version` 参数已经是这个带前缀的值，照抄进 `common.Migration{Version: version}`
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

## 硬约束：注册要赶在启动钩子之前

三个注册入口——`AppRouters`、`sdk.Runtime.SetAppRouters`、`migration.ForApp`——
都必须在 `cmd/api/server.go` 的 `runStartupHooks()` 执行之前调用完。

`init()` 是最省事的位置：Go 规范保证包级变量初始化与 `init()` 在 `main()` 之前
**单 goroutine 顺序执行**，注册期天然没有并发写。但它不是唯一合法位置——
在 `run()` 之类早于启动钩子的地方注册同样成立。这条规则约束的是**顺序**，
不是你写在哪个函数里。

`sdk.Runtime.SetAppRouters` 的准确语义以 core 为准：

> [go-admin-core `docs/contract.md`](https://github.com/go-admin-team/go-admin-core/blob/main/docs/contract.md)

那份文档写明了注册类与资源类的划分、封闭时刻、护栏边界（**只覆盖同步 panic，
你自己 `go func()` 出去的 panic 框架够不着**）、以及配置热更新会在运行期
重新执行 setup 回调这件事。

主仓这边只补三条它管不着的：

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
3. **`migration.ForApp` 是主仓的东西**，core 不认识它，上面那份文档不覆盖它。
   它的约束仍然是"注册要在迁移调度循环跑起来之前"，实践上就是 `init()`。

---

## 边界由 CI 守着

`common/`、`core/` 不得 import `app/`，这条由 `tools/checksilent` 的
`contract-import-boundary` 检查固化，`make checksilent` 在 CI 里跑，违反即失败
（测试文件同样算 —— 一个删掉 `app/admin` 的 fork 也应该能跑 `go test ./...`）。

靠人工评审列契约面会漏。上面那两处反向依赖里，第二处就是评审没发现、
靠机器全量扫描才找出来的。

`tools/checksilent` 还检查另外五类"不出声的失败"，写模块时值得先看一眼
`go run ./tools/checksilent -h`。
