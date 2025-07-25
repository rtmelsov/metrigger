package agent

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/metrics"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	err  error
	pkey *rsa.PublicKey
)

func Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	// created pem key files for asymmetric encrypting
	met := make(chan models.MetricsCollectorData)
	var PollCount float64
	logger := config.GetAgentConfig().GetLogger()

	cr := config.GetAgentConfig().GetCryptoKey()
	if cr != "" {
		pkey, err = helpers.LoadPublicKey(cr)
		if err != nil {
			logger.Error("Error to load a public key from environment's variable ", zap.String("error", err.Error()))
			return
		}
	}

	go metrics.CollectMetrics(PollCount, met)
	prettyJSON, _ := json.MarshalIndent(config.GetAgentConfig().Address(), "", "  ")
	logger.Info("started",
		zap.String("agent flags", string(prettyJSON)),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	)
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
					err := Worker(metricList, pkey)
					if err != nil {
						logger.Error("error while sending request to server", zap.String("error", err.Error()))
					}
				}()
			}
		}
	}
}
