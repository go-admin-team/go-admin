package cronjob

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go-admin/global"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/ws"
	"time"
)

var timeFormat = "2006-01-02 15:04:05"
var retryCount = 3

// 初始化
func JobSetup() {

	fmt.Println(time.Now().Format(timeFormat), " [INFO] Job Starting...")

	global.GADMCron = newWithSeconds()

	sysjob := models.SysJob{}
	joblist, err := sysjob.GetList()
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] Job init error", err)
	}
	if len(joblist) == 0 {
		fmt.Println(time.Now().Format(timeFormat), " [INFO] Job total:0")
	}

	_, err = sysjob.RemoveAllEntryID()
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] Job remove entry_id error", err)
	}

	for i := 0; i < len(joblist); i++ {
		AddJob(joblist[i])

	}

	// 其中任务
	global.GADMCron.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] Job start success.")
	// 关闭任务
	defer global.GADMCron.Stop()
	select {}
}

// 添加任务
func AddJob(job models.SysJob) error {
	h := HttpJob{Url: job.InvokeTarget, JobId: job.JobId, Name: job.JobName}
	// 添加定时任务
	id, err := global.GADMCron.AddJob(job.CronExpression, h)

	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] Job AddJob error", err)
		return err
	}
	job.EntryId = int(id)
	h.EntryId = job.EntryId
	_, err = job.Update(job.JobId)
	fmt.Println(time.Now().Format(timeFormat), " [INFO] Job AddJob success ,info:", job)
	return err
}

// 任务类型 http
type HttpJob struct {
	Url     string
	Name    string
	JobId   int
	EntryId int
}

//http 任务接口
func (h HttpJob) Run() {

	startTime := time.Now()
	var count = 0
	/* 循环 */
LOOP:
	if count < retryCount {
		/* 跳过迭代 */
		str, err := pkg.Get(h.Url)
		if err != nil {
			// 如果失败暂停一段时间重试
			fmt.Println(time.Now().Format(timeFormat), " [ERROR] mission failed! ", err)
			fmt.Printf(time.Now().Format(timeFormat)+" [INFO] Retry after the task fails %d seconds! %s \n", time.Duration(count)*time.Second, str)
			time.Sleep(time.Duration(count) * time.Second)
			goto LOOP
		}
		count = count + 1
		//fmt.Printf("a的值为 : %s \n", str)
	}

	fmt.Println(h.Url)
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	str := time.Now().Format(timeFormat) + " [INFO] Job " + string(h.EntryId) + "exec success , spend :" + latencyTime.String()
	ws.SendAll(str)
	global.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] Job ", h, "exec success , spend :", latencyTime)
}

// 调度停止
func Stop(cron *cron.Cron) chan bool {
	ch := make(chan bool)
	go func() {
		cron.Stop()
		ch <- true
	}()
	return ch
}

//移除Job
func Remove(e *cron.Cron, entryID int) {
	ch := make(chan bool)
	job := models.SysJob{}
	go func() {

		e.Remove(cron.EntryID(entryID))
		_, _ = job.RemoveEntryID(entryID)
		fmt.Println(time.Now().Format(timeFormat), " [INFO] Job Remove success ,info entryID :", entryID)
		ch <- true
	}()
}

// newWithSeconds returns a Cron with the seconds field enabled.
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
