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

## Configuration details

1. Configuration file description
```yml
settings:
  application:  
    # Project launch environment         
    env: dev  
    # When env: demo, prompts for request operations other than GET
    envmsg: "谢谢您的参与，但为了大家更好的体验，所以本次提交就算了吧！" 
    # Host IP or domain name, default 0.0.0.0
    host: 0.0.0.0 
    # Whether to initialize the database structure and basic data; true: required; false: not required
    isinit: false  
    # JWT encrypted string
    jwtsecret: 123abc  
    # log storage path
    logpath: temp/logs/log.log   
    # application name
    name: go-admin   
    # application port
    port: 8000   
    readtimeout: 1   
    writertimeout: 2 
  database:
    # database name
    database: dbname 
    # database type
    dbtype: mysql    
    # database host
    host: 127.0.0.1  
    # database  password
    password: password  
    # database port
    port: 3306       
    # database username
    username: root   
  redis:
    # redis addresss
    addr: 0.0.0.0:6379 
    # db 
    db: 0   
    # password            
    password: password  
    # read timeout
    readtimeout: 50     
```

2. file path  go-admin/config/settings.yml


## 📦 evelopment


First start instructions

```bash
# Get the code
git clone https://github.com/wenjianzhang/go-admin.git

# Enter working path
cd ./go-admin

# Build the project
go build

# Change setting
vi ./config/setting.yml (Note: Change isinit and database connection)

# 1. Database information in the configuration file
# Note: the corresponding configuration data under settings.database
# 2. Confirm database initialization parameters
# Note: If this is the first time settings.application.isinit is set, please set the current value to true, the system will automatically initialize the database structure and basic data information;
# 3. Confirm the log path


# Start the project or debug with the IDE
./go-admin

# See also instructions in WIKI
```


Document generation
```bash
swag init  
```

If there is no `swag` command go get installed
```bash
go get -u github.com/swaggo/swag/cmd/swag
```


Cross compilation
```bash
env GOOS=windows GOARCH=amd64 go build main.go

# or

env GOOS=linux GOARCH=amd64 go build main.go
```


## 🔗 Online Demo
> admin  /  123456

Demo address：[http://www.zhangwj.com](http://www.zhangwj.com/#/login)


## 🤝 Open source projects used
[gin](https://github.com/gin-gonic/gin)

[casbin](https://github.com/casbin/casbin)

[spf13/viper](https://github.com/spf13/viper)

[gorm](https://github.com/jinzhu/gorm)

[gin-swagger](https://github.com/swaggo/gin-swagger)

[jwt-go](https://github.com/dgrijalva/jwt-go)

[vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)

[ruoyi-vue](https://gitee.com/y_project/RuoYi-Vue)



## Version

#### 2020-03-15 New Features and Optimization

1. Add user avatar upload
2. Add user password modification
3. Operation log page adjustment
4. Optimize captcha background color

I saw a lot of friends who experience the wrong verification code, so I adjusted the contrast for everyone to experience!


## 🤝 Thanks
[chengxiao](https://github.com/chengxiao)


## License

[MIT](https://github.com/wenjianzhang/go-admin/blob/master/LICENSE.md)

Copyright (c) 2020 wenjianzhang

[中文]qq technical exchange group: 74520518
