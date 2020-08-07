package sd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"go-admin/tools/app"
	"net/http"
	"runtime"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// 健康状况
// @Summary 健康状况 HealthCheck shows OK as the ping-pong result.
// @Description 健康状况
// @Accept text/html
// @Produce text/html
// @Success 200 {string} string "OK"
// @Router /sd/health [get]
// @BasePath
func HealthCheck(c *gin.Context) {
	app.OK(c, "", "OK")
}

// @Summary 服务器硬盘使用量
// @Description 服务器硬盘使用量 DiskCheck checks the disk usage.
// @Accept  text/html
// @Produce text/html
// @Success 200 {string} string "OK - Free space: 16321MB (15GB) / 51200MB (50GB) | Used: 31%"
// @Failure 500 {string} string "CRITICAL"
// @Failure 429 {string} string "WARNING"
// @Router /sd/disk [get]
// @BasePath
func DiskCheck(c *gin.Context) {
	u, _ := disk.Usage("/")

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusOK
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.String(status, "\n"+message)
}

// @Summary OS
// @Description Os
// @Accept  text/html
// @Produce text/html
// @Success 200 {string} string ""
// @Router /sd/os [get]
// @BasePath
func OSCheck(c *gin.Context) {
	status := http.StatusOK
	app.Custum(c, gin.H{
		"code":         200,
		"status":       status,
		"goOs":         runtime.GOOS,
		"compiler":     runtime.Compiler,
		"numCpu":       runtime.NumCPU(),
		"version":      runtime.Version(),
		"numGoroutine": runtime.NumGoroutine(),
	})
}

// @Summary CPU 使用量
// @Description CPU 使用量 DiskCheck checks the disk usage.
// @Accept  text/html
// @Produce text/html
// @Success 200 {string} string ""
// @Router /sd/cpu [get]
// @BasePath
func CPUCheck(c *gin.Context) {
	cores, _ := cpu.Counts(false)

	cpus, err := cpu.Percent(time.Duration(200)*time.Millisecond, true)
	if err == nil {
		for i, c := range cpus {
			fmt.Printf("cpu%d : %f%%\n", i, c)
		}
	}

	a, _ := load.Avg()
	l1 := a.Load1
	l5 := a.Load5
	l15 := a.Load15

	status := http.StatusOK
	text := "OK"

	if l5 >= float64(cores-1) {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}
	app.Custum(c, gin.H{
		"code":   200,
		"msg":    text,
		"status": status,
		"cores":  cores,
		"load1":  l1,
		"load5":  l5,
		"load15": l15,
	})
}

// @Summary 内存使用量
// @Description 内存使用量 RAMCheck checks the disk usage.
// @Accept  text/html
// @Produce text/html
// @Success 200 {string} string ""
// @Router /sd/ram [get]
// @BasePath
func RAMCheck(c *gin.Context) {
	u, _ := mem.VirtualMemory()

	usedMB := int(u.Used) / MB
	totalMB := int(u.Total) / MB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	app.Custum(c, gin.H{
		"code":        200,
		"msg":         text,
		"status":      status,
		"used":        usedMB,
		"total":       totalMB,
		"usedPercent": usedPercent,
	})
}
