package status

import (
	"AIOPrivacyBot/utils"
	"fmt"
	"log"
	"runtime"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func getStatus() string {
	// 获取主机信息
	hostStat, err := host.Info()
	if err != nil {
		log.Printf("Error getting host info: %v", err)
		return "Error getting host info"
	}
	uptime := time.Duration(hostStat.Uptime) * time.Second

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Printf("Error getting CPU info: %v", err)
		return "Error getting CPU info"
	}
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU percent: %v", err)
		return "Error getting CPU percent"
	}

	// 获取内存信息
	memStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting virtual memory: %v", err)
		return "Error getting virtual memory"
	}
	swapStat, err := mem.SwapMemory()
	if err != nil {
		log.Printf("Error getting swap memory: %v", err)
		return "Error getting swap memory"
	}

	// 获取负载信息
	loadStat, err := load.Avg()
	if err != nil {
		log.Printf("Error getting load average: %v", err)
		return "Error getting load average"
	}

	// 系统信息格式化
	systemInfo := fmt.Sprintf(`
<b>系统信息:</b>
• <b>系统:</b> %s %s (%s)
• <b>运行时间:</b> %s
• <b>CPU:</b> %s
• <b>CPU 核心数:</b> %d
• <b>CPU 占用:</b> %.2f%%
• <b>内存占用:</b> %.2f MB / %.2f MB
• <b>交换内存:</b> %.2f MB / %.2f MB
• <b>1/5/15分钟负载:</b> %.2f / %.2f / %.2f
• <b>Go 版本:</b> %s
• <b>Goroutine 数量:</b> %d
`,
		hostStat.Platform, hostStat.PlatformVersion, hostStat.KernelVersion,
		uptime.String(),
		cpuInfo[0].ModelName,
		runtime.NumCPU(),
		cpuPercent[0],
		float64(memStat.Used)/1024/1024, float64(memStat.Total)/1024/1024,
		float64(swapStat.Used)/1024/1024, float64(swapStat.Total)/1024/1024,
		loadStat.Load1, loadStat.Load5, loadStat.Load15,
		runtime.Version(),
		runtime.NumGoroutine(),
	)

	return systemInfo
}

func HandleStatusCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	status := getStatus()
	err := utils.SendMessage(message.Chat.ID, status, message.MessageID, bot)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
