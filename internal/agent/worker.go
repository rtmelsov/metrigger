package agent

import (
	"context"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	pb "github.com/rtmelsov/metrigger/proto"
	"google.golang.org/grpc/metadata"
)

// Worker функция для отправки POST запроса
func Worker(metricList []*pb.Metric, c pb.MetricsServiceClient) error {
	// get logger method
	logger := config.GetAgentConfig().GetLogger()
	t := config.GetAgentConfig().TrustedSubnet()

	md := metadata.New(map[string]string{"X-Real-IP": t})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// Wrap a list of them:
	data, err := c.AddMetrics(ctx, &pb.AddMetricsRequest{
		Metrics: metricList,
	})
	logger.Info("get info")
	if err != nil {
		logger.Info("ERROR")
		return err
	}

	fmt.Printf("INFO: %v\r\n", data.Metrics)

	logger.Info("requested")
	return nil
}
