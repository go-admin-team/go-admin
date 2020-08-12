<p align="center">
  <img width="320" src="https://gitee.com/mydearzwj/image/raw/master/img/go-admin.svg">
</p>


<p align="center">
  <a href="https://github.com/wenjianzhang/go-admin">
    <img src="https://github.com/wenjianzhang/go-admin/workflows/build/badge.svg" alt="go-admin">
  </a>
  <a href="https://github.com/wenjianzhang/go-admin">
    <img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="license">
  </a>
    <a href="http://doc.zhangwj.com/go-admin-site/donate/">
    <img src="https://img.shields.io/badge/%24-donate-ff69b4.svg" alt="donate">
  </a>
</p>


  [English](https://github.com/wenjianzhang/go-admin/blob/master/README.en.md) | 简体中文


##### 基于Gin + Vue + Element UI的前后端分离权限管理系统 

系统初始化极度简单，只需要配置文件中，修改数据库连接，系统启动后会自动初始化数据库信息以及必须的基础数据

[在线文档国际](https://wenjianzhang.github.io/go-admin-site)

[在线文档国内](http://mydearzwj.gitee.io/go-admin-site/)

[前端项目](https://github.com/wenjianzhang/go-admin-ui)

[视频教程](https://space.bilibili.com/565616721/channel/detail?cid=125737)

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

- 多命令模式

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
11. 代码生成：根据数据表结构生成对应的增删改查相对应业务，全部可视化编程，基本业务可以0代码实现。
12. 表单构建：自定义页面样式，拖拉拽实现页面布局。
13. 服务监控：查看一些服务器的基本信息。

## 准备工作

你需要在本地安装 [go] [gin] [node](http://nodejs.org/) 和 [git](https://git-scm.com/) 

同时配套了系列教程包含视频和文档，如何从下载完成到熟练使用，强烈建议大家先看完这些教程再来实践本项目！！！

### 轻松实现go-admin写出第一个应用 - 文档教程

[步骤一 - 基础内容介绍](http://doc.zhangwj.com/go-admin-site/guide/intro/tutorial01.html)

[步骤二 - 实际应用 - 编写增删改查](http://doc.zhangwj.com/go-admin-site/guide/intro/tutorial02.html) 

### 手把手教你从入门到放弃 - 视频教程 

[如何启动go-admin](https://www.bilibili.com/video/BV1z5411x7JG)

[使用生成工具轻松实现业务](https://www.bilibili.com/video/BV1Dg4y1i79D)

[v1.1.0版本代码生成工具-释放双手](https://www.bilibili.com/video/BV1N54y1i71P) [进阶]

[多命令启动方式讲解以及IDE配置](https://www.bilibili.com/video/BV1Fg4y1q7ph)

[go-admin菜单的配置说明](https://www.bilibili.com/video/BV1Wp4y1D715) [必看]

[如何配置菜单信息以及接口信息](https://www.bilibili.com/video/BV1zv411B7nG) [必看]

[go-admin权限配置使用说明](https://www.bilibili.com/video/BV1rt4y197d3) [必看]

[go-admin数据权限使用说明](https://www.bilibili.com/video/BV1LK4y1s71e) [必看]


**如有问题请先看上述使用文档和文章，若不能满足，欢迎 issue 和 pr ，视频教程和文档持续更新中**

## 🗞 系统架构

<p align="center">
  <img  src="https://gitee.com/mydearzwj/image/raw/d9f59ea603e3c8a3977491a1bfa8f122e1a80824/img/go-admin-system.png" width="936px" height="491px">
</p>

## 📦 本地开发

### 开发目录创建

```bash

# 创建开发目录
mkdir goadmin
cd goadmin
```

### 获取代码

> 重点注意：两个项目必须放在同一文件夹下；

```bash
# 获取后端代码
git clone https://github.com/wenjianzhang/go-admin.git

# 获取前端代码
git clone https://github.com/wenjianzhang/go-admin-ui.git

```


### 启动说明

#### 服务端启动说明

```bash
# 进入 go-admin 后端项目
cd ./go-admin

# 编译项目
go build

# 修改配置 
# 文件路径  go-admin/config/settings.yml
vi ./config/setting.yml 

# 1. 配置文件中修改数据库信息 
# 注意: settings.database 下对应的配置数据
# 2. 确认log路径
```

#### 初始化数据库，以及服务启动
```
# 首次配置需要初始化数据库资源信息
./go-admin init -c config/settings.yml -m dev


# 启动项目，也可以用IDE进行调试
./go-admin server -c config/settings.yml -p 8000 -m dev

```

#### 使用docker 编译启动

```shell
# 编译镜像
docker build -t go-admin .

# 启动容器，第一个go-admin是容器名字，第二个go-admin是镜像名称
docker run --name go-admin -p 8000:8000 -d go-admin
```



#### 文档生成

```bash
swag init  

# 如果没有swag命令 go get安装一下即可
go get -u github.com/swaggo/swag/cmd/swag
```

#### 交叉编译
```bash
env GOOS=windows GOARCH=amd64 go build main.go

# or

env GOOS=linux GOARCH=amd64 go build main.go
```

### UI交互端启动说明

```bash
# 安装依赖
npm install

# 建议不要直接使用 cnpm 安装依赖，会有各种诡异的 bug。可以通过如下操作解决 npm 下载速度慢的问题
npm install --registry=https://registry.npm.taobao.org

# 启动服务
npm run dev
```

## 🎬 在线体验
> admin  /  123456

演示地址：[http://www.zhangwj.com](http://www.zhangwj.com/#/login)


## 📨 互动

<table>
  <tr>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/wx.png" width="180px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/qq.png" width="200px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/qq2.png" width="200px"></td>
  </tr>
  <tr>
    <td>微信</td>
    <td>此群已满</td>
    <td><a target="_blank" href="https://shang.qq.com/wpa/qunwpa?idkey=0f2bf59f5f2edec6a4550c364242c0641f870aa328e468c4ee4b7dbfb392627b"><img border="0" src="https://pub.idqqimg.com/wpa/images/group.png" alt="go-admin技术交流乙号" title="go-admin技术交流乙号"></a></td>
  </tr>
</table>
  

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

## 🤟 打赏

> 如果你觉得这个项目帮助到了你，你可以帮作者买一杯果汁表示鼓励 :tropical_drink:


<img class="no-margin" src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/pay.png"  height="200px" >

## ❤️ 赞助者

> 有部分是微信名称

zhuqiyun LLL狐 星星之火 cjj770 Sam 唐*i 晓聪 aLong *渊 海马 魏镇坪 + 111 *哥 我的宇哥哥 *声 *节

## 🔑 License

[MIT](https://github.com/wenjianzhang/go-admin/blob/master/LICENSE.md)

Copyright (c) 2020 wenjianzhang

