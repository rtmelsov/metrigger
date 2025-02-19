package agent

import (
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/metrics"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func Run() {
	met := make(chan models.MetricsCollector)
	var PollCount float64
	logger := config.GetAgentStorage().GetLogger()
	go metrics.CollectMetrics(PollCount, met)
	prettyJSON, _ := json.MarshalIndent(config.AgentFlags, "", "  ")
	logger.Info("started", zap.String("agent flags", string(prettyJSON)))
	for {
		var metricList []*models.Metrics
		time.Sleep(time.Duration(config.AgentFlags.ReportInterval) * time.Second)
		for k, b := range <-met {
			counter := RequestToServer("counter", k, 0, 1)
			gauge := RequestToServer("gauge", k, b, 0)
			metricList = append(metricList, counter, gauge)
		}
		data, err := json.Marshal(&metricList)
		logger := config.GetAgentStorage().GetLogger()
		if err != nil {
			logger.Panic("Error to Marshal JSON", zap.String("error", err.Error()))
			return
		}

		reqBody, err := helpers.CompressData(data)
		if err != nil {
			logger.Panic("Error to Marshal JSON", zap.String("error", err.Error()))
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
		respMetric, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Panic("error while to read", zap.String("error", err.Error()))
			return
		}

		logger.Info("respMetric", zap.String("metrc", string(respMetric)))

		err = resp.Body.Close()
		if err != nil {
			logger.Panic("3 Request to services", zap.String("error", err.Error()))
			return
		}
		logger.Info("requested")
	}
}

func RequestToServer(t string, key string, value float64, counter int64) *models.Metrics {
	var metric *models.Metrics

	if t == "counter" {
		metric = &models.Metrics{
			MType: t,
			ID:    key,
			Delta: &counter,
		}
	} else {
		metric = &models.Metrics{
			MType: t,
			ID:    key,
			Value: &value,
		}
	}

	return metric
}
