package metrics

import (
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/models"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// CollectMetrics - сбор метрик из runtime и gopsutil
func CollectMetrics(PollCount float64, m chan models.MetricsCollector) {
	for {
		time.Sleep(time.Duration(config.AgentFlags.PollInterval) * time.Second)

		// Сбор метрик runtime
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		var met = models.MetricsCollector{}
		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)

		RandomValue := r.Float64()
		PollCount++
		met["PollCount"] = PollCount
		met["RandomValue"] = RandomValue
		met["Alloc"] = float64(memStats.Alloc)
		met["HeapAlloc"] = float64(memStats.HeapAlloc)
		met["Sys"] = float64(memStats.Sys)
		met["NumGC"] = float64(memStats.NumGC)

		// 🔹 Сбор метрик памяти (TotalMemory, FreeMemory)
		v, _ := mem.VirtualMemory()
		met["TotalMemory"] = float64(v.Total)
		met["FreeMemory"] = float64(v.Free)

		// 🔹 Сбор загрузки CPU (CPUutilization1, CPUutilization2, ...)
		cpuUtilization, _ := cpu.Percent(time.Second, true)
		for i, val := range cpuUtilization {
			met["CPUutilization"+strconv.Itoa(i+1)] = val
		}

		// Отправка метрик в канал
		m <- met
	}
}
