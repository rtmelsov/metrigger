package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"runtime"
	"time"
)

type metrics map[string]float64

func Run(ctx context.Context) {
	met := make(chan metrics, 1) // Use a buffered channel to prevent deadlocks.

	logger := storage.GetMemStorage().GetLogger()
	logger.Info("Agent started")

	// Goroutine to collect runtime memory statistics
	go func(ctx context.Context, m chan metrics) {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Metrics collection stopped")
				return
			default:
				time.Sleep(time.Duration(config.AgentFlags.PollInterval) * time.Second)
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)

				met := metrics{
					"Alloc":         float64(memStats.Alloc),
					"BuckHashSys":   float64(memStats.BuckHashSys),
					"Frees":         float64(memStats.Frees),
					"GCCPUFraction": float64(memStats.GCCPUFraction),
					"GCSys":         float64(memStats.GCSys),
					"HeapAlloc":     float64(memStats.HeapAlloc),
					"HeapIdle":      float64(memStats.HeapIdle),
					"HeapInuse":     float64(memStats.HeapInuse),
					"HeapObjects":   float64(memStats.HeapObjects),
					"HeapReleased":  float64(memStats.HeapReleased),
					"HeapSys":       float64(memStats.HeapSys),
					"LastGC":        float64(memStats.LastGC),
					"Lookups":       float64(memStats.Lookups),
					"MCacheInuse":   float64(memStats.MCacheInuse),
					"MCacheSys":     float64(memStats.MCacheSys),
					"MSpanInuse":    float64(memStats.MSpanInuse),
					"MSpanSys":      float64(memStats.MSpanSys),
					"Mallocs":       float64(memStats.Mallocs),
					"NextGC":        float64(memStats.NextGC),
					"NumForcedGC":   float64(memStats.NumForcedGC),
					"NumGC":         float64(memStats.NumGC),
					"OtherSys":      float64(memStats.OtherSys),
					"PauseTotalNs":  float64(memStats.PauseTotalNs),
					"StackInuse":    float64(memStats.StackInuse),
					"StackSys":      float64(memStats.StackSys),
					"Sys":           float64(memStats.Sys),
					"TotalAlloc":    float64(memStats.TotalAlloc),
				}
				m <- met
			}
		}
	}(ctx, met)

	// Report metrics to the server
	for {
		select {
		case <-ctx.Done():
			logger.Info("Agent reporting stopped")
			return
		default:
			time.Sleep(time.Duration(config.AgentFlags.ReportInterval) * time.Second)
			for k, b := range <-met {
				RequestToServer("counter", k, 0, 1)
				RequestToServer("gauge", k, b, 0)
			}
			logger.Info("Metrics reported")
		}
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
	logger := storage.GetMemStorage().GetLogger()
	if err != nil {
		logger.Error("Failed to marshal JSON", zap.Error(err))
		return
	}

	requestBody := bytes.NewReader(data)
	url := fmt.Sprintf("http://%s/update/", config.AgentFlags.Addr)

	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		logger.Error("Failed to create HTTP request", zap.Error(err))
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Failed to send request to server", zap.Error(err))
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Error("Failed to close response body", zap.Error(closeErr))
		}
	}()
}
