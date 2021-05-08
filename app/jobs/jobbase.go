package jobs

import (
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	models2 "go-admin/app/jobs/models"
	"gorm.io/gorm"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/cronjob"
)

var timeFormat = "2006-01-02 15:04:05"
var retryCount = 3

var jobList map[string]JobsExec
var lock sync.Mutex

type JobCore struct {
	InvokeTarget   string
	Name           string
	JobId          int
	EntryId        int
	CronExpression string
	Args           string
}

// 任务类型 http
type HttpJob struct {
	JobCore
}

type ExecJob struct {
	JobCore
}

func (e *ExecJob) Run() {
	startTime := time.Now()
	var obj = jobList[e.InvokeTarget]
	if obj == nil {
		log.Warn("[Job] ExecJob Run job nil")
		return
	}
	err := CallExec(obj.(JobsExec), e.Args)
	if err != nil {
		// 如果失败暂停一段时间重试
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] mission failed! ", err)
	}
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	//TODO: 待完善部分
	//str := time.Now().Format(timeFormat) + " [INFO] JobCore " + string(e.EntryId) + "exec success , spend :" + latencyTime.String()
	//ws.SendAll(str)
	log.Info("[Job] JobCore %s exec success , spend :%v", e.Name, latencyTime)
	return
}

//http 任务接口
func (h *HttpJob) Run() {

	startTime := time.Now()
	var count = 0
	var err error
	var str string
	/* 循环 */
LOOP:
	if count < retryCount {
		/* 跳过迭代 */
		str, err = pkg.Get(h.InvokeTarget)
		if err != nil {
			// 如果失败暂停一段时间重试
			fmt.Println(time.Now().Format(timeFormat), " [ERROR] mission failed! ", err)
			fmt.Printf(time.Now().Format(timeFormat)+" [INFO] Retry after the task fails %d seconds! %s \n", (count+1)*5, str)
			time.Sleep(time.Duration(count+1) * 5 * time.Second)
			count = count + 1
			goto LOOP
		}
	}
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	//TODO: 待完善部分

	log.Infof("[Job] JobCore %s exec success , spend :%v", h.Name, latencyTime)
	return
}

// 初始化
func Setup(dbs map[string]*gorm.DB) {

	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Starting...")

	for k, db := range dbs {
		sdk.Runtime.SetCrontab(k, cronjob.NewWithSeconds())
		setup(k, db)
	}
}

func setup(key string, db *gorm.DB) {
	crontab := sdk.Runtime.GetCrontabKey(key)
	sysJob := models2.SysJob{}
	jobList := make([]models2.SysJob, 0)
	err := sysJob.GetList(db, &jobList)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore init error", err)
	}
	if len(jobList) == 0 {
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore total:0")
	}

	_, err = sysJob.RemoveAllEntryID(db)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore remove entry_id error", err)
	}

	for i := 0; i < len(jobList); i++ {
		if jobList[i].JobType == 1 {
			j := &HttpJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName

			sysJob.EntryId, err = AddJob(crontab, j)
		} else if jobList[i].JobType == 2 {
			j := &ExecJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName
			j.Args = jobList[i].Args
			sysJob.EntryId, err = AddJob(crontab, j)
		}
		err = sysJob.Update(db, jobList[i].JobId)
	}

	// 其中任务
	crontab.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore start success.")
	// 关闭任务
	defer crontab.Stop()
	select {}
}

// 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(c *cron.Cron, job Job) (int, error) {
	if job == nil {
		fmt.Println("unknown")
		return 0, nil
	}
	return job.addJob(c)
}

func (h *HttpJob) addJob(c *cron.Cron) (int, error) {
	id, err := c.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

func (h *ExecJob) addJob(c *cron.Cron) (int, error) {
	id, err := c.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

// 移除任务
func Remove(c *cron.Cron, entryID int) chan bool {
	ch := make(chan bool)
	go func() {
		c.Remove(cron.EntryID(entryID))
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Remove success ,info entryID :", entryID)
		ch <- true
	}()
	return ch
}

// 任务停止
//func Stop() chan bool {
//	ch := make(chan bool)
//	go func() {
//		global.GADMCron.Stop()
//		ch <- true
//	}()
//	return ch
//}
