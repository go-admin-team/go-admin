package monitor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"go-admin/tools"
	"go-admin/tools/app"
	"runtime"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// @Summary 系统信息
// @Description 获取JSON
// @Tags 系统信息
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/settings/serverInfo [get]
func ServerInfo(c *gin.Context) {

	osDic := make(map[string]interface{}, 0)
	osDic["goOs"] = runtime.GOOS
	osDic["arch"] = runtime.GOARCH
	osDic["mem"] = runtime.MemProfileRate
	osDic["compiler"] = runtime.Compiler
	osDic["version"] = runtime.Version()
	osDic["numGoroutine"] = runtime.NumGoroutine()
	osDic["ip"] = tools.GetLocaHonst()
	osDic["projectDir"] = tools.GetCurrentPath()

	dis, _ := disk.Usage("/")
	diskTotalGB := int(dis.Total) / GB
	diskFreeGB := int(dis.Free) / GB
	diskDic := make(map[string]interface{}, 0)
	diskDic["total"] = diskTotalGB
	diskDic["free"] = diskFreeGB

	mem, _ := mem.VirtualMemory()
	memUsedMB := int(mem.Used) / GB
	memTotalMB := int(mem.Total) / GB
	memFreeMB := int(mem.Free) / GB
	memUsedPercent := int(mem.UsedPercent)
	memDic := make(map[string]interface{}, 0)
	memDic["total"] = memTotalMB
	memDic["used"] = memUsedMB
	memDic["free"] = memFreeMB
	memDic["usage"] = memUsedPercent

	cpuDic := make(map[string]interface{}, 0)
	cpuDic["cpuInfo"],_ = cpu.Info()
	percent,_ := cpu.Percent(0,false)
	cpuDic["Percent"] = fmt.Sprintf("%.2f",percent[0])
	cpuDic["cpuNum"], _ = cpu.Counts(false)

	app.Custum(c, gin.H{
		"code": 200,
		"os":   osDic,
		"mem":  memDic,
		"cpu":  cpuDic,
		"disk": diskDic,
	})
}
