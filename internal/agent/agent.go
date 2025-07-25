package agent

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/interfaces"
	"github.com/rtmelsov/metrigger/internal/metrics"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Agent type for declare agent (client) global variables
type Agent struct {
	pkey   *rsa.PublicKey
	logger *zap.Logger
	config interfaces.AgentActionsI
}

func NewAgent() *Agent {
	a := &Agent{
		logger: config.GetAgentConfig().GetLogger(),
		config: config.GetAgentConfig(),
	}
	if cr := a.config.GetCryptoKey(); cr != "" {
		key, err := helpers.LoadPublicKey(cr)
		if err != nil {
			a.logger.Error("failed to load pulbic key", zap.String("error", err.Error()))
		} else {
			a.pkey = key
		}
	}

	return a
}

func (a *Agent) Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	jobs := make(chan []*models.Metrics, 100)
	var wg sync.WaitGroup

	workerCount := a.config.RateLimit()
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			defer wg.Done()
			for task := range jobs {
				if err := Worker(task, a.pkey); err != nil {
					a.logger.Error("worker error", zap.Int("id", id), zap.Error(err))
				}
			}
		}(i)
	}

	// created pem key files for asymmetric encrypting
	met := make(chan models.MetricsCollectorData)

	go metrics.CollectMetrics(0, met)

	t := time.NewTicker(time.Duration(a.config.ReportInterval()) * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("waiting...")
			close(jobs)
			wg.Wait()
			return
		case <-t.C:
			a.logger.Info("tick")
			var metricData = <-met
			metricList := make([]*models.Metrics, 0, metricData.Length)
			for k, b := range *metricData.Metrics {
				metricList = append(
					metricList,
					RequestToServer("counter", k, 0, 1),
					RequestToServer("gauge", k, b, 0),
				)
			}

			jobs <- metricList
		}
	}
}
