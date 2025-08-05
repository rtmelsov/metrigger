// Package middleware
package middleware

import (
	"context"
	"github.com/rtmelsov/metrigger/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TrustedSubnet(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	storage.GetMemStorage().GetLogger().Info("trusted subnet check")

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			trustedSubnet := values[0]
			if trustedSubnet != storage.ServerFlags.TrustedSubnet {
				return nil, status.Errorf(codes.Aborted, "exptected subnet %v, but we got %v", storage.ServerFlags.TrustedSubnet, trustedSubnet)
			}
		}
	}

	storage.GetMemStorage().GetLogger().Info("we don't have any trusted subnet")

	return handler(ctx, req)
}
