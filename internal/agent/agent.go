package agent

import (
	"context"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/metrics"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	met := make(chan models.MetricsCollectorData)
	var PollCount float64
	logger := config.GetAgentConfig().GetLogger()
	go metrics.CollectMetrics(PollCount, met)
	t := time.NewTicker(time.Duration(config.GetAgentConfig().ReportInterval()) * time.Second)

	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			fmt.Println("waiting...")
			wg.Wait()
			return
		case <-t.C:
			logger.Info("tick")
			var metricData = <-met
			metricList := make([]*models.Metrics, 0, metricData.Length)
			for k, b := range *metricData.Metrics {
				counter := RequestToServer("counter", k, 0, 1)
				gauge := RequestToServer("gauge", k, b, 0)
				metricList = append(metricList, counter, gauge)
			}

			rl := config.GetAgentConfig().RateLimit()
			wg.Add(rl)
			for i := 0; i < rl; i++ {
				go func() {
					defer wg.Done()
					err := Worker(metricList)
					if err != nil {
						logger.Error("error while sending request to server", zap.String("error", err.Error()))
					}
				}()
			}
		}
	}
}
