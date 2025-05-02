package agent

import (
	"encoding/json"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/metrics"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"time"
)

func Run() {
	met := make(chan models.MetricsCollectorData)
	var PollCount float64
	logger := config.GetAgentConfig().GetLogger()
	go metrics.CollectMetrics(PollCount, met)
	prettyJSON, _ := json.MarshalIndent(config.GetAgentConfig().Address(), "", "  ")
	logger.Info("started",
		zap.String("agent flags", string(prettyJSON)),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	)
	t := time.NewTicker(time.Duration(config.GetAgentConfig().ReportInterval()) * time.Second)
	for range t.C {
		logger.Info("tick")
		var metricData = <-met
		metricList := make([]*models.Metrics, 0, metricData.Length)
		for k, b := range *metricData.Metrics {
			counter := RequestToServer("counter", k, 0, 1)
			gauge := RequestToServer("gauge", k, b, 0)
			metricList = append(metricList, counter, gauge)
		}

		for i := 0; i < config.GetAgentConfig().RateLimit(); i++ {
			go func() {
				err := Worker(metricList)
				if err != nil {
					logger.Error("error while sending request to server", zap.String("error", err.Error()))
				}
			}()
		}
	}
}
