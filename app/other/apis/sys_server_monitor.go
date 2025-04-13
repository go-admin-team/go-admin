package apis

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var excludeNetInterfaces = []string{
	"lo", "tun", "docker", "veth", "br-", "vmbr", "vnet", "kube",
}

type ServerMonitor struct {
	api.Api
}

// GetHourDiffer 获取相差时间
func GetHourDiffer(startTime, endTime string) int64 {
	t1, err1 := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err2 := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err1 != nil || err2 != nil || !t1.Before(t2) {
		return 0
	}
	return (t2.Unix() - t1.Unix()) / 3600
}

// ServerInfo 获取系统信息
func (e ServerMonitor) ServerInfo(c *gin.Context) {
	e.Context = c

	osInfo := getOSInfo()
	memInfo := getMemoryInfo()
	swapInfo := getSwapInfo()
	cpuInfo := getCPUInfo()
	diskInfo := getDiskInfo()
	netInfo := getNetworkInfo()

	bootTime, _ := host.BootTime()
	cachedBootTime := time.Unix(int64(bootTime), 0)

	e.Custom(gin.H{
		"code":     200,
		"os":       osInfo,
		"mem":      memInfo,
		"cpu":      cpuInfo,
		"disk":     diskInfo,
		"net":      netInfo,
		"swap":     swapInfo,
		"location": "Aliyun",
		"bootTime": GetHourDiffer(cachedBootTime.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05")),
	})
}

func getOSInfo() map[string]interface{} {
	sysInfo, _ := host.Info()
	return map[string]interface{}{
		"goOs":         runtime.GOOS,
		"arch":         runtime.GOARCH,
		"mem":          runtime.MemProfileRate,
		"compiler":     runtime.Compiler,
		"version":      runtime.Version(),
		"numGoroutine": runtime.NumGoroutine(),
		"ip":           pkg.GetLocalHost(),
		"projectDir":   pkg.GetCurrentPath(),
		"hostName":     sysInfo.Hostname,
		"time":         time.Now().Format("2006-01-02 15:04:05"),
	}
}

func getMemoryInfo() map[string]interface{} {
	memory, _ := mem.VirtualMemory()
	return map[string]interface{}{
		"used":    memory.Used / MB,
		"total":   memory.Total / MB,
		"percent": pkg.Round(memory.UsedPercent, 2),
	}
}

func getSwapInfo() map[string]interface{} {
	memory, _ := mem.VirtualMemory()
	return map[string]interface{}{
		"used":  memory.SwapTotal - memory.SwapFree,
		"total": memory.SwapTotal,
	}
}

func getCPUInfo() map[string]interface{} {
	cpuInfo, _ := cpu.Info()
	percent, _ := cpu.Percent(0, false)
	cpuNum, _ := cpu.Counts(false)
	return map[string]interface{}{
		"cpuInfo": cpuInfo,
		"percent": pkg.Round(percent[0], 2),
		"cpuNum":  cpuNum,
	}
}

func getDiskInfo() map[string]interface{} {
	var diskTotal, diskUsed, diskUsedPercent float64
	diskList := make([]disk.UsageStat, 0)

	diskInfo, err := disk.Partitions(true)
	if err == nil {
		for _, p := range diskInfo {
			diskDetail, err := disk.Usage(p.Mountpoint)
			if err == nil {
				diskDetail.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskDetail.UsedPercent), 64)
				diskDetail.Total /= MB
				diskDetail.Used /= MB
				diskDetail.Free /= MB
				diskList = append(diskList, *diskDetail)
			}
		}
	}

	d, _ := disk.Usage("/")
	diskTotal = float64(d.Total / GB)
	diskUsed = float64(d.Used / GB)
	diskUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", d.UsedPercent), 64)

	return map[string]interface{}{
		"total":   diskTotal,
		"used":    diskUsed,
		"percent": diskUsedPercent,
	}
}

func getNetworkInfo() map[string]interface{} {
	netInSpeed, netOutSpeed := trackNetworkSpeed()
	return map[string]interface{}{
		"in":  pkg.Round(float64(netInSpeed/KB), 2),
		"out": pkg.Round(float64(netOutSpeed/KB), 2),
	}
}

func trackNetworkSpeed() (uint64, uint64) {
	var netInSpeed, netOutSpeed, netInTransfer, netOutTransfer, lastUpdateNetStats uint64
	nc, err := net.IOCounters(true)
	if err == nil {
		for _, v := range nc {
			if isListContainsStr(excludeNetInterfaces, v.Name) {
				continue
			}
			netInTransfer += v.BytesRecv
			netOutTransfer += v.BytesSent
		}
		now := uint64(time.Now().Unix())
		diff := now - lastUpdateNetStats
		if diff > 0 {
			netInSpeed = (netInTransfer - netInTransfer) / diff
			netOutSpeed = (netOutTransfer - netOutTransfer) / diff
		}
		lastUpdateNetStats = now
	}
	return netInSpeed, netOutSpeed
}

func isListContainsStr(list []string, str string) bool {
	for _, item := range list {
		if strings.Contains(str, item) {
			return true
		}
	}
	return false
}
