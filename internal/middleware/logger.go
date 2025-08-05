package middleware

import (
	"context"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	storage.GetMemStorage().GetLogger().Info("Incoming request...")

	resp, err := handler(ctx, req)

	storage.GetMemStorage().GetLogger().Info(
		"Finished!",
		zap.String("full_method", info.FullMethod),
		zap.Duration("duration", time.Since(start)),
	)
	if err != nil {
		st, _ := status.FromError(err)
		storage.GetMemStorage().GetLogger().Error(st.Message())
	}

	return resp, nil
}
