package agent

import (
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type metrics map[string]float64

func Run() {
	met := make(chan metrics)
	var PollCount float64
	logger := config.GetAgentStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(config.AgentFlags, "", "  ")
	logger.Info("started", zap.String("agent flags", string(prettyJSON)))

	go func(m chan metrics) {
		for {
			time.Sleep(time.Duration(config.AgentFlags.PollInterval) * time.Second)
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			var met = metrics{}
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
	}(met)
	for {
		time.Sleep(time.Duration(config.AgentFlags.ReportInterval) * time.Second)
		for k, b := range <-met {
			RequestToServer("counter", k, 0, 1)
			RequestToServer("gauge", k, b, 0)
		}

		logger.Info("requested")
	}
}

func RequestToServer(t string, key string, value float64, counter int64) {
	var metric models.Metrics

	if t == "counter" {
		metric = models.Metrics{
			MType: t,
			ID:    key,
			Delta: &counter,
		}
	} else {
		metric = models.Metrics{
			MType: t,
			ID:    key,
			Value: &value,
		}
	}
	data, err := json.Marshal(metric)
	logger := config.GetAgentStorage().GetLogger()
	if err != nil {
		logger.Panic("Error to Marshal SSON", zap.String("error", err.Error()))
		return
	}

	reqBody, err := helpers.CompressData(data)
	if err != nil {
		logger.Panic("Error to Marshal SSON", zap.String("error", err.Error()))
		return
	}

	url := fmt.Sprintf("http://%s/update/", config.AgentFlags.Addr)

	req, err := http.NewRequest("POST", url, reqBody)

	if err != nil {
		logger.Panic("1 Request to services", zap.String("error", err.Error()))
		return
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logger.Panic("2 Request to services", zap.String("error", err.Error()))
		return
	}
	err = resp.Body.Close()
	if err != nil {
		logger.Panic("3 Request to services", zap.String("error", err.Error()))
		return
	}

}
