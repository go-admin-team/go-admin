package sys_job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"go-admin/app/admin/service"
	"go-admin/common/apis"
	"go-admin/common/dto"
)

type SysJob struct {
	apis.Api
}

// RemoveJobForService 调用service实现
func (e SysJob) RemoveJobForService(c *gin.Context) {
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}
	var v dto.GeneralDelDto
	err = c.BindUri(&v)
	if err != nil {
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	s := service.SysJob{}
	s.Log = log
	s.Orm = db
	s.Cron = sdk.Runtime.GetCrontabKey(c.Request.Host)
	err = s.RemoveJob(&v)
	if err != nil {
		log.Errorf("RemoveJob error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(nil, s.Msg)
}

// StartJobForService 启动job service实现
func (e SysJob) StartJobForService(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}
	var v dto.GeneralGetDto
	err = c.BindUri(&v)
	if err != nil {
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	s := service.SysJob{}
	s.Orm = db
	s.Log = log
	s.Cron = sdk.Runtime.GetCrontabKey(c.Request.Host)
	err = s.StartJob(&v)
	if err != nil {
		log.Errorf("GetCrontabKey error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(nil, s.Msg)
}
