#!/bin/bash
echo "go build"
go mod tidy
go build -o go-admin main.go
chmod +x ./go-admin
echo "kill go-admin service"
killall go-admin # kill go-admin service
nohup ./go-admin server -c=config/settings.dev.yml >> access.log 2>&1 & #后台启动服务将日志写入access.log文件
echo "run go-admin success"
ps -aux | grep go-admin
