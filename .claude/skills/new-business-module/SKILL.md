---
name: new-business-module
description: Scaffold a new single-table CRUD business module end to end — migration, Actions-mode model/dto/router, and the sys_menu/sys_api/casbin seed data that makes it show up in the UI with working permissions. Use when the user wants to add a new business table/module to go-admin, not for cross-table or non-CRUD business logic.
---

# 新增业务模块

给一张新的业务表配齐"能跑、能看见、能授权"的完整闭环：迁移 → 后端代码 → 菜单与权限种子数据。
只适用于单表增删改查；跨表事务、外部调用、复杂校验等超出这个范围（见下方"何时不适用"）。

开始前先读 `AGENTS.md`（分层边界、通用 Action 使用前提、命名规则）和 `app/demo/` 下的全部文件——
这是可编译、有测试、CI 会跑的参照物，本文与它冲突时以它为准。

## 何时不适用

业务超出单表 CRUD（跨表事务、外部服务调用、复杂校验）时，不要用这个 skill 硬套——
改成手写 Handler + Service，参照 `app/admin/apis/sys_post.go` 及其 Service，遵守
`AGENTS.md` 的分层约束（Api 不碰 Orm，Service 不碰 `gin.Context`，一律用 `e.Orm`）。

## 步骤

### 1. 确认表结构

表结构需符合命名规范：`sys_`/业务前缀 + 下划线（如 `tb_article`）。核对字段是否已有
`created_at`/`updated_at`/`deleted_at` 这类约定字段。

### 2. 写数据库迁移

放在 `cmd/migrate/migration/version/` 目录（**不是** `version-local/` —— 后者在
`.gitignore` 中，提交时会被忽略，`git status` 也看不到）。

- 文件名前 13 位是时间戳版本号
- 已执行过的迁移文件不可修改；需要修正时新增一个迁移
- 包名为 `version`

### 3. 生成 model / dto / router 三个文件（Actions 模式）

不要手写 Api 与 Service。使用 `common/actions` 的通用 Action，一个模块只需
model、dto、router 三个文件，完整写法照抄 `app/demo/` 的结构。

**关键正确性要求**（这三条是实际出问题最多的地方）：

- Model 实现 `models.ActiveRecord`（`Generate` / `GetId` / `TableName`），
  `TableName()` 必须显式声明——GORM 配置了 `SingularTable`，不会自动推导
- **`Generate()` 必须返回副本，不要就地返回**——Action 在并发请求间复用实例，
  就地返回会导致请求之间串数据；这个问题单人测试时几乎不出现，上线后才暴露
- 完成后确认 `cmd/api/` 中已用 `_` 导入新包，否则路由不会被注册

### 4. 写菜单、接口与权限种子数据

这一步最容易被漏掉——代码能编译、接口能测通，但界面上看不到菜单、点了按钮说
没权限，往往就是漏了这一步。结构参照 `cmd/migrate/migration/version/1786700001000_demo_menu.go`
——它是可运行、幂等（用 `upsert`，重复跑不会报错）的真实例子。

:::danger
**但不要照抄它的 import。** 那个文件用的是 `cmd/migrate/migration/models`，
只因为它的版本号排在软删除转换（`1786700003000`）之前才是安全的。

**你新写的迁移版本号在转换之后，必须改用 `app/` 下的运行时模型**
（`app/admin/models.SysApi`、`SysMenu`），否则第一条 insert 就会
`NOT NULL constraint failed: sys_api.deleted_at`。
`TestPostConversionMigrationsAvoidFrozenSeedModels` 会拦住这个错误。
:::

一个模块要在界面上可用，需要四类数据，缺一样都不行：

| 表 | 作用 |
|---|---|
| `sys_api` | 后端路由登记，Casbin 据此判定权限 |
| `sys_menu` | 侧边栏菜单（目录用 `M`、菜单用 `C`、按钮用 `F`） |
| `sys_menu_api_rule` | 菜单与接口的多对多关联，角色保存时据此生成策略 |
| `casbin_rule` | 实际生效的权限策略（**不是** `sys_casbin_rule`，那张表的唯一索引在 MySQL 下会超长，不要迁移它） |

必须核对的两处一致性——**错了不会报错，只会在界面上表现为"看不到/点不动"**：

- `sys_menu.menu_name` 必须与前端组件的 `defineOptions({ name: 'XxxManage' })` 一致，
  否则 `keep-alive` 缓存静默失效
- 按钮级 `sys_menu.permission`（格式 `模块:资源:操作`）必须与前端
  `v-permisaction="['模块:资源:操作']"` 完全一致，否则按钮权限判断静默失效

### 5. 收尾检查

| 检查项 | 出错后果 |
| --- | --- |
| `Generate()` 是否返回副本 | 并发请求之间串数据 |
| 是否使用 `e.Orm` 而非全局 DB | 多租户下拿到错误的数据库连接 |
| `TableName()` 是否显式声明 | GORM 不会自动推导 |
| 迁移文件是否放在 `version/` | 放进 `version-local/` 会被忽略，别人拉代码看不到 |
| `sys_menu.menu_name` 是否与前端组件 `name` 一致 | keep-alive 缓存静默失效 |
| `sys_menu.permission` 是否与前端 `v-permisaction` 一致 | 按钮权限静默失效 |

跑一遍 `go run -tags sqlite3 . migrate -c config/settings.sqlite.yml` 验证迁移可执行，
再 `go run -tags sqlite3 . server -c config/settings.sqlite.yml` 启动服务，用 admin
账号登录确认新菜单和按钮权限都出现了。

如果前端页面还没生成，下一步用 go-admin-ui 仓库里的 `new-list-page` skill——两边靠
`sys_menu.permission` / `v-permisaction` 这个字符串对齐。

> 不要把 `config/settings.yml` 的真实内容贴给 AI 工具——`database.source` 含数据库
> 账号密码，`jwt.secret` 泄露后可被用来伪造任意用户的 token。
