package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/middleware"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"strings"
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

// Webhook функция для распределения адресов для определения методов
func Webhook() chi.Router {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)

	privateKey, err := helpers.LoadPrivateKey(storage.ServerFlags.CryptoRate)
	if err != nil {
		panic(privateKey)
	}
	r.Use(middleware.GzipParser)
	r.Use(middleware.JwtParser)
	r.Use(middleware.CryptoParser(privateKey))
	r.Route("/", func(r chi.Router) {
		r.Get("/", MetricsListHandler)
		r.Get("/ping", PingDBHandler)
		r.Route("/updates", MetricsUpdateListHandler)
		r.Route("/update", MetricsUpdateHandler)
		r.Route("/value", MetricsValueHandler)
	})

	return r
}
