package agent

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"context"
	"crypto/rsa"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/interfaces"
	"github.com/rtmelsov/metrigger/internal/metrics"
	"github.com/rtmelsov/metrigger/internal/models"
	pb "github.com/rtmelsov/metrigger/proto"
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
	return &Agent{
		logger: config.GetAgentConfig().GetLogger(),
		config: config.GetAgentConfig(),
	}
}

func (a *Agent) Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		a.logger.Error("worker error", zap.Error(err))
	}
	defer conn.Close()

	c := pb.NewMetricsServiceClient(conn)

	jobs := make(chan []*pb.Metric, 100)
	var wg sync.WaitGroup

	workerCount := a.config.RateLimit()
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			defer wg.Done()
			for task := range jobs {
				if err := Worker(task, c); err != nil {
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
			metricList := make([]*pb.Metric, 0, metricData.Length)
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
