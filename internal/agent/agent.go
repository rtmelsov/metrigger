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
	met := make(chan models.MetricsCollector)
	var PollCount float64
	logger := config.GetAgentStorage().GetLogger()
	go metrics.CollectMetrics(PollCount, met)
	prettyJSON, _ := json.MarshalIndent(config.AgentFlags, "", "  ")
	logger.Info("started",
		zap.String("agent flags", string(prettyJSON)),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	)
	for {
		time.Sleep(time.Duration(config.AgentFlags.ReportInterval) * time.Second)
		var metricList []*models.Metrics
		for k, b := range <-met {
			counter := RequestToServer("counter", k, 0, 1)
			gauge := RequestToServer("gauge", k, b, 0)
			metricList = append(metricList, counter, gauge)
		}

		for i := 0; i < config.AgentFlags.RateLimit; i++ {
			go worker(metricList)
		}
	}
}
