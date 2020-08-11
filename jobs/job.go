package jobs

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go-admin/global"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/cronjob"
	"go-admin/pkg/ws"
	"reflect"
	"time"
)

var timeFormat = "2006-01-02 15:04:05"
var retryCount = 3

type Job struct {
	InvokeTarget   string
	Name           string
	JobId          int
	EntryId        int
	CronExpression string
}

// 任务类型 http
type HttpJob struct {
	Job
}

type ExecJob struct {
	Job
}

func (e ExecJob) Run() {
	startTime := time.Now()

	if result := callReflect(&EXEC{}, e.InvokeTarget); result != nil {
		fmt.Printf("callReflectMethod ShowMs %s \n", result[0].String())
	} else {
		fmt.Println("callReflectMethod ShowMs didn't run ")
	}

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)
	str := time.Now().Format(timeFormat) + " [INFO] Job " + string(e.EntryId) + "exec success , spend :" + latencyTime.String()
	ws.SendAll(str)
	global.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] Job ", e, "exec success , spend :", latencyTime)
}

//http 任务接口
func (h HttpJob) Run() {

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
	str := time.Now().Format(timeFormat) + " [INFO] Job " + string(h.EntryId) + "exec success , spend :" + latencyTime.String()
	ws.SendAll(str)
	global.JobLogger.Info(time.Now().Format(timeFormat), " [INFO] Job ", h, "exec success , spend :", latencyTime)
}

// 初始化
func Setup() {

	fmt.Println(time.Now().Format(timeFormat), " [INFO] Job Starting...")

	global.GADMCron = cronjob.NewWithSeconds()

	sysJob := models.SysJob{}
	jobList, err := sysJob.GetList()
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] Job init error", err)
	}
	if len(jobList) == 0 {
		fmt.Println(time.Now().Format(timeFormat), " [INFO] Job total:0")
	}

	_, err = sysJob.RemoveAllEntryID()
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] Job remove entry_id error", err)
	}

	for i := 0; i < len(jobList); i++ {
		if jobList[i].JobType == 1 {
			j:=HttpJob{}
			j.InvokeTarget=jobList[i].InvokeTarget
			j.CronExpression=jobList[i].CronExpression
			j.JobId=jobList[i].JobId
			j.Name=jobList[i].JobName

			sysJob.EntryId, err = AddJob(j)
		} else if jobList[i].JobType == 2 {
			j:=ExecJob{}
			j.InvokeTarget=jobList[i].InvokeTarget
			j.CronExpression=jobList[i].CronExpression
			j.JobId=jobList[i].JobId
			j.Name=jobList[i].JobName
			sysJob.EntryId, err = AddJob(j)
		}
		_, err = sysJob.Update(jobList[i].JobId)
	}

	// 其中任务
	global.GADMCron.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] Job start success.")
	// 关闭任务
	defer global.GADMCron.Stop()
	select {}
}

// 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(job interface{}) (int, error) {
	switch job.(type) {
	case HttpJob:
		op, ok := job.(HttpJob)
		if ok {
			return op.addJob()
		}
	case ExecJob:
		op, ok := job.(ExecJob)
		if ok {
			return op.addJob()
		}
	default:
		fmt.Println("unknown")
		return 0, nil
	}
	fmt.Println("job error")
	return 0, nil
}

func (h HttpJob) addJob() (int, error) {
	id, err := global.GADMCron.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] Job AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

func (h ExecJob) addJob() (int, error) {
	id, err := global.GADMCron.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] Job AddJob error", err)
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
		fmt.Println(time.Now().Format(timeFormat), " [INFO] Job Remove success ,info entryID :", entryID)
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

func callReflect(any interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	if v := reflect.ValueOf(any).MethodByName(name); v.String() == "<invalid Value>" {
		return nil
	} else {
		return v.Call(inputs)
	}

}
