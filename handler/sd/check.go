package sd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"go-admin/utils"
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
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
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
// @Success 200 {string} string "OK - Free space: 455MB (0GB) / 8192MB (8GB) | Used: 5%"
// @Failure 500 {string} string "CRITICAL"
// @Failure 429 {string} string "WARNING"
// @Router /sd/os [get]
// @BasePath
func OSCheck(c *gin.Context) {
	fmt.Println(runtime.GOOS)     //取得操作系统版本
	fmt.Println(runtime.Compiler) //取得编译器名称-gc
	fmt.Println(runtime.NumCPU()) //取得CPU多核的数目
	fmt.Println(runtime.Version())

	status := http.StatusOK
	c.String(status, runtime.GOOS+"/"+runtime.GOARCH+"/"+utils.IntToString(runtime.NumCPU()))
}

// @Summary CPU 使用量
// @Description CPU 使用量 DiskCheck checks the disk usage.
// @Accept  text/html
// @Produce text/html
// @Success 200 {string} string "OK - Load average: 2.39, 2.13, 1.97 | Cores: 2"
// @Failure 500 {string} string "CRITICAL"
// @Failure 429 {string} string "WARNING"
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

	fmt.Println(c)

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

	message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)
	c.String(status, "\n"+message)
}

// @Summary 内存使用量
// @Description 内存使用量 RAMCheck checks the disk usage.
// @Accept  text/html
// @Produce text/html
// @Success 200 {string} string "OK - Free space: 455MB (0GB) / 8192MB (8GB) | Used: 5%"
// @Failure 500 {string} string "CRITICAL"
// @Failure 429 {string} string "WARNING"
// @Router /sd/ram [get]
// @BasePath
func RAMCheck(c *gin.Context) {
	u, _ := mem.VirtualMemory()

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
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

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.String(status, "\n"+message)
}
