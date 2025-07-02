package handlers

import (
	"github.com/rtmelsov/metrigger/internal/middleware"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetMetricData(r *http.Request) (string, string) {
	logger := storage.GetMemStorage().GetLogger()
	paths := strings.Split(r.URL.String(), "/")
	logger.Debug("paths: %v\r\n", zap.String("paths", strings.Join(paths, ", ")))
	var metname, metval string
	if len(paths) > 3 {
		metname = paths[3]
	}
	if len(paths) > 4 {
		metval = paths[4]
	}

	return metname, metval
}

func Webhook() chi.Router {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Use(middleware.GzipParser)
	r.Use(middleware.JwtParser)
	r.Route("/", func(r chi.Router) {
		r.Get("/", MerticsListHandler)
		r.Get("/ping", PingDBHandler)
		r.Route("/updates", MetricsUpdateListHandler)
		r.Route("/update", MetricsUpdateHandler)
		r.Route("/value", MetricsValueHandler)
	})

	return r
}
