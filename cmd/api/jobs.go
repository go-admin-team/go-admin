package api

import "go-admin/app/jobs/router"

func init() {
	//注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)

	// test
}
