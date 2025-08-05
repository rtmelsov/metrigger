package handlers

import (
	"context"
	pb "github.com/rtmelsov/metrigger/proto"

	"github.com/rtmelsov/metrigger/internal/services"
	"github.com/rtmelsov/metrigger/internal/storage"
)

func (w *WebhookServer) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse
	var err error

	response.Metric, err = services.UpdateMetric(in.Metric)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WebhookServer) AddMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var response pb.AddMetricsResponse
	var err error
	storage.GetMemStorage().GetLogger().Info("request func: JSONUpdateList")

	if storage.ServerFlags.DataBaseDsn != "" {
		storage.GetMemStorage().GetLogger().Info("in db")
		response.Metrics, err = UpdateMetrics(in.Metrics)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	dataList := make([]*pb.Metric, 0)
	for _, v := range in.Metrics {
		metric, err := services.UpdateMetric(v)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, metric)
	}
	response.Metrics = dataList
	return &response, nil
}
