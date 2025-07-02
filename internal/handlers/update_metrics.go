package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/services"
)

func MetricsUpdateHandler(r chi.Router) {
	UpdateRequests := map[string]func(string, string) error{
		"counter": services.MetricsCounterSet,
		"gauge":   services.MetricsGaugeSet,
	}
	r.Post("/", JSONUpdate)
	for k := range UpdateRequests {
		r.Route(fmt.Sprintf("/%s", k), func(r chi.Router) {
			r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
				if fn, exist := UpdateRequests[k]; exist {
					metName, metVal := GetMetricData(r)
					if err := SetMeticsUpdate(metName, metVal, fn); err != nil {
						http.Error(w, err.Text, err.StatusCode)
						return
					}
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte("success"))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					http.Error(w, "Can't find parameters", http.StatusBadRequest)
					return
				}
			})
		})
	}

	r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unknown type", http.StatusBadRequest)
	})
}

func SetMeticsUpdate(metName, metVal string, fn func(string, string) error) *models.ErrorType {
	if metName == "" || metVal == "" {
		return &models.ErrorType{
			Text: "can't find parameters", StatusCode: http.StatusNotFound,
		}
	}
	if err := fn(metName, metVal); err != nil {
		return &models.ErrorType{
			Text: "can't find parameters", StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}
