#  go-admin  


##### 基于Gin + Vue + Element UI的前后端分离权限管理系统


## ✨ 特性

- 遵循 RESTful API 设计规范

- 基于 GIN WEB API 框架，提供了丰富的中间件支持（用户认证、跨域、访问日志、追踪ID等）

- 基于Casbin的 RBAC 访问控制模型

- JWT 认证

- 支持 Swagger 文档(基于swaggo)

- 基于 GORM 的数据库存储，可扩展多种类型数据库 

- 配置文件简单的模型映射，快速能够得到想要的配置

- TODO: 单元测试


## 🎁 内置

1.  用户管理：用户是系统操作者，该功能主要完成系统用户配置。
2.  部门管理：配置系统组织机构（公司、部门、小组），树结构展现支持数据权限。
3.  岗位管理：配置系统用户所属担任职务。
4.  菜单管理：配置系统菜单，操作权限，按钮权限标识等。
5.  角色管理：角色菜单权限分配、设置角色按机构进行数据范围权限划分。
6.  字典管理：对系统中经常使用的一些较为固定的数据进行维护。
7.  参数管理：对系统动态配置常用参数。
8.  操作日志：系统正常操作日志记录和查询；系统异常信息日志记录和查询。
9.  登录日志：系统登录日志记录查询包含登录异常。
10. 系统接口：根据业务代码自动生成相关的api接口文档。


## 📦 本地开发


step 1:
```bash
git clone https://e.coding.net/wenjianzhang/go-admin.git
```
step 2:
```bash
cd ./goadmin/src/goadmin
```
step 3:
```bash
go build
```

step 4:
```bash
vi ./config/setting.yml (更改isinit和数据库连接)
```
step 5:
```bash
./goadmin
```

文档生成
```bash
swag init  
```

如果没有swag命令 go get安装一下即可
```bash
go get -u github.com/swaggo/swag/cmd/swag
```

交叉编译
```bash
env GOOS=windows GOARCH=amd64 go build main.go
```
or
```bash
env GOOS=linux GOARCH=amd64 go build main.go
```

## 🔗 在线体验
> admin/123456

演示地址：http://www.zhangwj.com


## 🤝 使用的开源项目
[gin](https://github.com/gin-gonic/gin)

[casbin](https://github.com/casbin/casbin)

[spf13/viper](https://github.com/spf13/viper)

[gorm](https://github.com/jinzhu/gorm)

[gin-swagger](https://github.com/swaggo/gin-swagger)

[jwt-go](https://github.com/dgrijalva/jwt-go)

[vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)

[ruoyi-vue](https://gitee.com/y_project/RuoYi-Vue)



特别感谢
[chengxiao](https://github.com/chengxiao)
