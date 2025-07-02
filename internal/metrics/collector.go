// Package metrics предназначен для сбора метрик из runtime и gopsutil
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

// CollectMetric получает счетчик и добавляет другие метрики runtime в объект и возвращает мап со структурой и длину мап
func CollectMetric(PollCount *float64) (*models.MetricsCollector, int) {
	// Сбор метрик runtime
	var length int
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	var met = models.MetricsCollector{}
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	RandomValue := r.Float64()
	*PollCount++
	length++
	met["PollCount"] = *PollCount
	length++
	met["RandomValue"] = RandomValue
	length++
	met["Alloc"] = float64(memStats.Alloc)
	length++
	met["HeapAlloc"] = float64(memStats.HeapAlloc)
	length++
	met["Sys"] = float64(memStats.Sys)
	length++
	met["NumGC"] = float64(memStats.NumGC)
	length++

	// Сбор метрик памяти (TotalMemory, FreeMemory)
	v, _ := mem.VirtualMemory()
	met["TotalMemory"] = float64(v.Total)
	length++
	met["FreeMemory"] = float64(v.Free)
	length++

	// Сбор загрузки CPU (CPUutilization1, CPUutilization2, ...)
	cpuUtilization, _ := cpu.Percent(time.Second, true)
	for i, val := range cpuUtilization {
		met["CPUutilization"+strconv.Itoa(i+1)] = val
		length++
	}
	return &met, length
}

// CollectMetrics - сбор метрик из runtime и gopsutil
func CollectMetrics(PollCount float64, m chan models.MetricsCollectorData) {
	for {
		time.Sleep(time.Duration(config.GetAgentConfig().PollInterval()) * time.Second)
		met, length := CollectMetric(&PollCount)
		m <- models.MetricsCollectorData{
			Metrics: met,
			Length:  length,
		}
	}
}
