#  go-admin  

![build](https://github.com/wenjianzhang/go-admin/workflows/build/badge.svg)   ![license](https://img.shields.io/github/license/mashape/apistatus.svg) 
  
  [English](./README.en.md) | 简体中文
  

##### 基于Gin + Vue + Element UI的前后端分离权限管理系统 

系统初始化极度简单，只需要配置文件中，修改数据库连接，系统启动后会自动初始化数据库信息以及必须的基础数据

[在线文档](https://wenjianzhang.github.io/go-admin-site)

## ✨ 特性

- 遵循 RESTful API 设计规范

- 基于 GIN WEB API 框架，提供了丰富的中间件支持（用户认证、跨域、访问日志、追踪ID等）

- 基于Casbin的 RBAC 访问控制模型

- JWT 认证

- 支持 Swagger 文档(基于swaggo)

- 基于 GORM 的数据库存储，可扩展多种类型数据库 

- 配置文件简单的模型映射，快速能够得到想要的配置

- 代码生成工具

- 表单构建工具

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
11. 代码生成：根据数据表结构生成对应的增删改查相对应业务，全部可视化编程。
12. 表单构建：自定义页面样式，拖拉拽实现页面布局。

## ⚙ 配置详情

1. 配置文件说明
```yml
settings:
  application:  
    # 项目启动环境            
    env: dev  
    # 当 env:demo 时，GET以外的请求操作提示
    envmsg: "谢谢您的参与，但为了大家更好的体验，所以本次提交就算了吧！" 
    # 主机ip 或者域名，默认0.0.0.0
    host: 0.0.0.0 
    # 是否需要初始化数据库结构以及基本数据；true：需要；false：不需要 
    isinit: false  
    # JWT加密字符串
    jwtsecret: 123abc  
    # log存放路径
    logpath: temp/logs/log.log   
    # 服务名称
    name: go-admin   
    # 服务端口
    port: 8000   
    readtimeout: 1   
    writertimeout: 2 
  database:
    # 数据库名称
    database: dbname 
    # 数据库类型
    dbtype: mysql    
    # 数据库地址
    host: 127.0.0.1  
    # 数据库密码
    password: password  
    # 数据库端口
    port: 3306       
    # 数据库用户名
    username: root   
```

2. 文件路径  go-admin/config/settings.yml


## 📦 本地开发

首次启动说明

```bash
# 获取代码
git clone https://github.com/wenjianzhang/go-admin.git

# 进入工作路径
cd ./go-admin

# 编译项目
go build

# 修改配置
vi ./config/setting.yml (更改isinit和数据库连接)

# 1. 配置文件中修改数据库信息 
# 注意: settings.database 下对应的配置数据)
# 2. 确认数据库初始化参数 
# 注意: settings.application.isinit 如果是首次启动，请把当前值设置成true，系统会自动初始化数据库结构以及基本的数据信息；
# 3. 确认log路径


# 启动项目，也可以用IDE进行调试
./go-admin

```
[在线文档](https://wenjianzhang.github.io/go-admin-site)


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

# or

env GOOS=linux GOARCH=amd64 go build main.go
```


## 🔗 在线体验
> admin  /  123456

演示地址：[http://www.zhangwj.com](http://www.zhangwj.com/#/login)

## 📈 版本

### 2020-04-08 新功能及优化

1. 添加sqlite3的支持
1. 数据库字段格式统一
2. 用户新增bug修复
3. 修改数据初始化脚本
4. 验证码改为数字验证 
5. 删除redis暂时无用组件
6. 其他已知bug的修复

### 2020-04-01 新功能及优化

1. 代码生成器
2. 代码优化
3. 已知bug修复

#### 2020-03-15 新功能及优化

1. 添加用户头像上传
2. 添加用户密码修改
3. 操作日志页面调整
4. 优化验证码背景色

看到好多体验的朋友验证码错误，所以调整了对比度，方便大家体验！

## 📨 互动

<table>
  <tr>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/go-admin/master/demo/wx.png" width="180px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/go-admin/master/demo/qq.png" width="200px"></td>
  </tr>
  <tr>
    <td>微信</td>
    <td>QQ</td>
  </tr>
</table>
  
<a target="_blank" href="https://shang.qq.com/wpa/qunwpa?idkey=1affb445445bd442312fcad9a927007db74a0cd4380bbc08a6c97d2691744869"><img border="0" src="https://pub.idqqimg.com/wpa/images/group.png" alt="go-admin技术交流" title="go-admin技术交流"></a>


## 🤝 特别感谢
[chengxiao](https://github.com/chengxiao)
[gin](https://github.com/gin-gonic/gin)
[casbin](https://github.com/casbin/casbin)
[spf13/viper](https://github.com/spf13/viper)
[gorm](https://github.com/jinzhu/gorm)
[gin-swagger](https://github.com/swaggo/gin-swagger)
[jwt-go](https://github.com/dgrijalva/jwt-go)
[vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)
[ruoyi-vue](https://gitee.com/y_project/RuoYi-Vue)


## License

[MIT](https://github.com/wenjianzhang/go-admin/blob/master/LICENSE.md)

Copyright (c) 2020 wenjianzhang
