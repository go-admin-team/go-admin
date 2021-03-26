#  go-admin  

![build](https://github.com/wenjianzhang/go-admin/workflows/build/badge.svg)   ![license](https://img.shields.io/github/license/mashape/apistatus.svg) 


 English | [简体中文](./README.zh-CN.md)
  
##### Gin + Vue + Element UI based scaffolding for front and back separation management system


## ✨ Feature

- Follow RESTful API design specifications

- Provides rich middleware support based on GIN WEB API framework (user authentication, cross domain, access log, tracking ID, etc.)

- Casbin-based RBAC access control model

- JWT certification

- Support Swagger documentation (based on swaggo)

- GORM-based database storage that can expand many types of databases

- Simple model mapping of configuration files to quickly get the desired configuration
       
- Form builder
  
- Multi command mode

- TODO: unit test


## 🎁 Built-in functions

1.  User management: The user is the system operator. This function mainly completes the system user configuration.
2.  Department management: configure the system organization (company, department, group), and display the tree structure to support data permissions.
3.  Post management: Configure system users to hold positions.
4.  Menu management: configure system menus, operation permissions, button permission labels, etc.
5.  Role management: role menu permissions assignment, setting roles to divide data range permissions by organization.
6.  Dictionary management: to maintain some fixed data often used in the system.
7.  Parameter management: Dynamically configure common parameters for the system.
8.  Operation log: system normal operation log record and query; system exception information log record and query.
9.  Login log: The system login log record query contains login exceptions.
10. System interface: Automatically generate related api interface documents according to business code.
11. Multi command mode code generation: according to the data table structure to generate the corresponding business, all visual programming, basic business can be 0 code.
12. Form construction: Customize page style, drag and drop to achieve page layout.
13. Service monitoring: view the basic information of some servers.       

## 🗞 System architecture

<p align="center">
  <img  src="https://gitee.com/mydearzwj/image/raw/d9f59ea603e3c8a3977491a1bfa8f122e1a80824/img/go-admin-system.png" width="936px" height="491px">
</p>

## 📦 Local development

### Create development directory

```bash
mkdir goadmin
cd goadmin
```

### Get code

> Key points: two projects must be placed in the same folder;

```bash
# Get backend code
git clone https://github.com/go-admin-team/go-admin.git

# Get front end code
git clone https://github.com/go-admin-team/go-admin-ui.git

```


### Start up instructions

#### Server startup instructions

```bash
# Enter the go-admin backend project
cd ./go-admin

# Compile project
go build

# Configuration 
# file path  go-admin/config/settings.yml
vi ./config/setting.yml 

# 1. Modifying database information in configuration file 
# settings.database Corresponding configuration data under
# 2. Confirm log path
```

#### Initialize the database and start the service
```
# The first configuration needs to initialize the database resource information
./go-admin migrate -c config/settings.yml


# Start the project, you can also debug with IDE
./go-admin server -c config/settings.yml

```

#### Using docker to compile and start

```shell
# Compile image
docker build -t go-admin .

# Start the container. The first go admin is the container name and the second go admin is the image name
docker run --name go-admin -p 8000:8000 -v /config/settings.yml:/config/settings.yml -d go-admin

```



#### Document generation

```bash
go generate
```

#### Cross compiling
```bash
env GOOS=windows GOARCH=amd64 go build main.go

# or

env GOOS=linux GOARCH=amd64 go build main.go
```

### UI interactive terminal startup description

```bash
# Installation dependency
npm install

# It is recommended that you do not use cnpm to install dependencies directly. There will be all kinds of weird bugs. The problem of slow download speed of NPM can be solved by the following operations
npm install --registry=https://registry.npm.taobao.org

# Start the service
npm run dev
```

## 🎬 Online Demo
> admin  /  123456

Demo：[http://www.go-admin.dev](http://www.go-admin.dev/#/login)


## 📨 Interaction

<table>
  <tr>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/wx.png" width="180px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/qq.png" width="200px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/qq2.png" width="200px"></td>
  </tr>
  <tr>
    <td>WeChat</td>
    <td>This QQ group is full</td>
    <td><a target="_blank" href="https://shang.qq.com/wpa/qunwpa?idkey=0f2bf59f5f2edec6a4550c364242c0641f870aa328e468c4ee4b7dbfb392627b"><img border="0" src="https://pub.idqqimg.com/wpa/images/group.png" alt="go-admin技术交流乙号" title="go-admin技术交流乙号"></a></td>
  </tr>
</table>

## 💎 Member

<a href="https://github.com/wenjianzhang"> <img src="https://avatars.githubusercontent.com/u/3890175?s=460&u=20eac63daef81588fbac611da676b99859319251&v=4" width="80px"></a>
<a href="https://github.com/lwnmengjing"> <img src="https://avatars.githubusercontent.com/u/12806223?s=400&u=a89272dce50100b77b4c0d5c81c718bf78ebb580&v=4" width="80px"></a>
<a href="https://github.com/chengxiao"> <img src="https://avatars.githubusercontent.com/u/1379545?s=460&u=557da5503d0ac4a8628df6b4075b17853d5edcd9&v=4" width="80px"></a>
<a href="https://github.com/bing127"> <img src="https://avatars.githubusercontent.com/u/31166183?s=460&u=c085bff88df10bb7676c8c0351ba9dcd031d1fb3&v=4" width="80px"></a>



## JetBrains Open source certificate support

The go admin project has always been developed in the GoLand integrated development environment of JetBrains company, based on the * * free JetBrains open source license (s) * * genuine free license. I would like to express my thanks.

<a href="https://www.jetbrains.com/?from=kubeadm-ha" target="_blank"><img src="https://raw.githubusercontent.com/panjf2000/illustrations/master/jetbrains/jetbrains-variant-4.png" width="250" align="middle"/></a>




## 🤝 Open source projects used
[gin](https://github.com/gin-gonic/gin)

[casbin](https://github.com/casbin/casbin)

[spf13/viper](https://github.com/spf13/viper)

[gorm](https://github.com/jinzhu/gorm)

[gin-swagger](https://github.com/swaggo/gin-swagger)

[jwt-go](https://github.com/dgrijalva/jwt-go)

[vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)

[ruoyi-vue](https://gitee.com/y_project/RuoYi-Vue)




## 🤝 Thanks
2. [chengxiao](https://github.com/chengxiao)



## License

[MIT](https://github.com/wenjianzhang/go-admin/blob/master/LICENSE.md)

Copyright (c) 2020 wenjianzhang

[中文]qq technical exchange group: 74520518