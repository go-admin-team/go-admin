# ⚙ 配置详情

1. 配置文件说明
```yml
settings:
  application:  
    # 项目启动环境            
    mode: dev  # dev开发环境 test测试环境 prod线上环境；
    host: 0.0.0.0  # 主机ip 或者域名，默认0.0.0.0
    # 服务名称
    name: go-admin   
    # 服务端口
    port: 8000   
    readtimeout: 1   
    writertimeout: 2 
  log:
    # 日志文件存放路径
    dir: temp/logs
  jwt:
    # JWT加密字符串
    secret: go-admin
    # 过期时间单位：秒
    timeout: 3600
  database:
    # 数据库名称
    name: dbname 
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