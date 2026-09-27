package app

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

var (
	appName  string
	StartCmd = &cobra.Command{
		Use:     "app",
		Short:   "Create a new app",
		Long:    "Use when you need to create a new app",
		Example: "go-admin app -n admin",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&appName, "name", "n", "", "Start server with provided configuration file")
}

func run() {

	fmt.Println(`start init`)
	//1. 读取配置

	fmt.Println(`generate migration file`)
	if err := genFile(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func genFile() error {
	if appName == "" {
		return errors.New("arg `name` invalid ：name is empty")
	}
	path := "app/"
	appPath := path + appName
	err := utils.IsNotExistMkDir(appPath)
	if err != nil {
		return err
	}
	apiPath := appPath + "/apis/"
	err = utils.IsNotExistMkDir(apiPath)
	if err != nil {
		return err
	}
	modelsPath := appPath + "/models/"
	err = utils.IsNotExistMkDir(modelsPath)
	if err != nil {
		return err
	}
	routerPath := appPath + "/router/"
	err = utils.IsNotExistMkDir(routerPath)
	if err != nil {
		return err
	}
	servicePath := appPath + "/service/"
	err = utils.IsNotExistMkDir(servicePath)
	if err != nil {
		return err
	}
	dtoPath := appPath + "/service/dto/"
	err = utils.IsNotExistMkDir(dtoPath)
	if err != nil {
		return err
	}

	t1, err := template.ParseFiles("template/cmd_api.template")
	if err != nil {
		return err
	}
	m := map[string]string{}
	m["appName"] = appName
	var b1 bytes.Buffer
	if err = t1.Execute(&b1, m); err != nil {
		return err
	}
	t2, err := template.ParseFiles("template/router.template")
	if err != nil {
		return err
	}
	var b2 bytes.Buffer
	if err = t2.Execute(&b2, nil); err != nil {
		return err
	}
	if err = pkg.FileCreate(b1, "./cmd/api/"+appName+".go"); err != nil {
		return err
	}
	return pkg.FileCreate(b2, appPath+"/router/router.go")
}
