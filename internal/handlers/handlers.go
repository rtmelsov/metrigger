package handlers

import (
	pb "github.com/rtmelsov/metrigger/proto"

	"google.golang.org/grpc"
)

type WebhookServer struct {
	pb.UnimplementedMetricsServiceServer
}

// InitWebhook функция для распределения адресов для определения методов
func InitWebhook(s *grpc.Server) {
	pb.RegisterMetricsServiceServer(s, &WebhookServer{})
}
