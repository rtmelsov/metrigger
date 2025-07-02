package metrics

import (
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/models"
	"math/rand"
	"runtime"
	"time"
)

func CollectMetrics(PollCount float64, m chan models.MetricsCollector) {
	for {
		time.Sleep(time.Duration(config.AgentFlags.PollInterval) * time.Second)
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		var met = models.MetricsCollector{}
		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)

		// Генерация случайного целого числа
		RandomValue := r.Float64() // Случайное число от 0 до 99
		PollCount++
		met["PollCount"] = PollCount
		met["RandomValue"] = RandomValue
		met["Alloc"] = float64(memStats.Alloc)
		met["BuckHashSys"] = float64(memStats.BuckHashSys)
		met["Frees"] = float64(memStats.Frees)
		met["GCCPUFraction"] = float64(memStats.GCCPUFraction)
		met["GCSys"] = float64(memStats.GCSys)
		met["HeapAlloc"] = float64(memStats.HeapAlloc)
		met["HeapIdle"] = float64(memStats.HeapIdle)
		met["HeapInuse"] = float64(memStats.HeapInuse)
		met["HeapObjects"] = float64(memStats.HeapObjects)
		met["HeapReleased"] = float64(memStats.HeapReleased)
		met["HeapSys"] = float64(memStats.HeapSys)
		met["LastGC"] = float64(memStats.LastGC)
		met["Lookups"] = float64(memStats.Lookups)
		met["MCacheInuse"] = float64(memStats.MCacheInuse)
		met["MCacheSys"] = float64(memStats.MCacheSys)
		met["MSpanInuse"] = float64(memStats.MSpanInuse)
		met["MSpanSys"] = float64(memStats.MSpanSys)
		met["Mallocs"] = float64(memStats.Mallocs)
		met["NextGC"] = float64(memStats.NextGC)
		met["NumForcedGC"] = float64(memStats.NumForcedGC)
		met["NumGC"] = float64(memStats.NumGC)
		met["OtherSys"] = float64(memStats.OtherSys)
		met["PauseTotalNs"] = float64(memStats.PauseTotalNs)
		met["StackInuse"] = float64(memStats.StackInuse)
		met["StackSys"] = float64(memStats.StackSys)
		met["Sys"] = float64(memStats.Sys)
		met["TotalAlloc"] = float64(memStats.TotalAlloc)
		m <- met
	}
}
