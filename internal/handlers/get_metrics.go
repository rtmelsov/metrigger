package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/services"
	"github.com/rtmelsov/metrigger/internal/storage"
	"net/http"
)

func MetricsValueHandler(r chi.Router) {
	ValueRequests := map[string]func(name string) (*models.CounterMetric, *models.GaugeMetric, error){
		"counter": services.MetricsCounterGet,
		"gauge":   services.MetricsGaugeGet,
	}
	r.Post("/", JSONGet)
	for k := range ValueRequests {
		r.Route(fmt.Sprintf("/%s", k), func(r chi.Router) {
			r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				storage.GetMemStorage().GetLogger().Info("request func: get value handler")
				if fn, exist := ValueRequests[k]; exist {
					metName, extra := GetMetricData(r)
					counter, gauge, err := GetMetricsValue(metName, extra, fn)
					if err != nil {
						http.Error(w, err.Text, err.StatusCode)
						return
					}
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					switch k {
					case "counter":
						err := SendData(w, []byte(fmt.Sprintf("%v", counter.Value)))
						if err != nil {
							http.Error(w, "Can't send response", http.StatusInternalServerError)
						}
						return
					case "gauge":
						err := SendData(w, []byte(fmt.Sprintf("%v", gauge.Value)))
						if err != nil {
							http.Error(w, "Can't send response", http.StatusInternalServerError)
						}
						return
					default:
						http.Error(w, "Can't find parameters", http.StatusBadRequest)
						return
					}
				} else {
					http.Error(w, "Can't find parameters", http.StatusBadRequest)
					return
				}
			})
		})
	}
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unknown type", http.StatusBadRequest)
	})
}

func GetMetricsValue(
	metName, extra string,
	fn func(name string) (*models.CounterMetric, *models.GaugeMetric, error),
) (*models.CounterMetric, *models.GaugeMetric, *models.ErrorType) {
	if extra != "" {
		return nil, nil, &models.ErrorType{
			Text: "can't find parameters", StatusCode: http.StatusNotFound,
		}
	}
	counter, gauge, err := fn(metName)
	if err != nil {
		return nil, nil, &models.ErrorType{
			Text: "can't find parameters", StatusCode: http.StatusNotFound,
		}
	}

	return counter, gauge, nil
}
