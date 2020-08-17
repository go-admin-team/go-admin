package jobs

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go-admin/global"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/cronjob"
	"sync"
	"time"
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
		global.JobLogger.Warning(" ExecJob Run job nil", e)
		return
	}
	CallExec(obj.(JobsExec))
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	//TODO: 待完善部分
	//str := time.Now().Format(timeFormat) + " [INFO] JobCore " + string(e.EntryId) + "exec success , spend :" + latencyTime.String()
	//ws.SendAll(str)
	global.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] JobCore ", e, "exec success , spend :", latencyTime)
}

//http 任务接口
func (h *HttpJob) Run() {

	startTime := time.Now()
	var count = 0
	/* 循环 */
LOOP:
	if count < retryCount {
		/* 跳过迭代 */
		str, err := pkg.Get(h.InvokeTarget)
		if err != nil {
			// 如果失败暂停一段时间重试
			fmt.Println(time.Now().Format(timeFormat), " [ERROR] mission failed! ", err)
			fmt.Printf(time.Now().Format(timeFormat)+" [INFO] Retry after the task fails %d seconds! %s \n", time.Duration(count)*time.Second, str)
			time.Sleep(time.Duration(count) * time.Second)
			goto LOOP
		}
		count = count + 1
	}
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	//TODO: 待完善部分
	//str := time.Now().Format(timeFormat) + " [INFO] JobCore " + string(h.EntryId) + "exec success , spend :" + latencyTime.String()
	//ws.SendAll(str)
	global.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] JobCore ", h, "exec success , spend :", latencyTime)
}

// 初始化
func Setup() {

	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Starting...")

	global.GADMCron = cronjob.NewWithSeconds()

	sysJob := models.SysJob{}
	jobList, err := sysJob.GetList()
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore init error", err)
	}
	if len(jobList) == 0 {
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore total:0")
	}

	_, err = sysJob.RemoveAllEntryID()
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

			sysJob.EntryId, err = AddJob(j)
		} else if jobList[i].JobType == 2 {
			j := &ExecJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName
			sysJob.EntryId, err = AddJob(j)
		}
		_, err = sysJob.Update(jobList[i].JobId)
	}

	// 其中任务
	global.GADMCron.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore start success.")
	// 关闭任务
	defer global.GADMCron.Stop()
	select {}
}

// 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(job Job) (int, error) {
	if job == nil {
		fmt.Println("unknown")
		return 0, nil
	}
	return job.addJob()
}

func (h *HttpJob) addJob() (int, error) {
	id, err := global.GADMCron.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

func (h *ExecJob) addJob() (int, error) {
	id, err := global.GADMCron.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

// 移除任务
func Remove(entryID int) chan bool {
	ch := make(chan bool)
	go func() {
		global.GADMCron.Remove(cron.EntryID(entryID))
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Remove success ,info entryID :", entryID)
		ch <- true
	}()
	return ch
}

// 任务停止
func Stop() chan bool {
	ch := make(chan bool)
	go func() {
		global.GADMCron.Stop()
		ch <- true
	}()
	return ch
}
