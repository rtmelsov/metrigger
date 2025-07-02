package middleware

import (
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r)
		duration := time.Since(start)
		logger := storage.GetMemStorage().GetLogger()
		logger.Info("URL data:",
			zap.String("url", uri),
			zap.String("method", method),
			zap.Float64("duration", duration.Seconds()))
	})
}
