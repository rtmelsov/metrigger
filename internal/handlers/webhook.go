package handlers

import (
	"fmt"
	"github.com/rtmelsov/metrigger/internal/middleware"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rtmelsov/metrigger/internal/constants"
	"github.com/rtmelsov/metrigger/internal/server"
)

type ErrorType struct {
	text       string
	statusCode int
}

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

func SetMeticsUpdate(w http.ResponseWriter, r *http.Request, fn func(string, string) error) *ErrorType {
	metName, metVal := GetMetricData(r)
	if metName == "" || metVal == "" {
		return &ErrorType{
			text: "can't find parameters", statusCode: http.StatusNotFound,
		}
	}
	if err := fn(metName, metVal); err != nil {
		return &ErrorType{
			text: "can't find parameters", statusCode: http.StatusBadRequest,
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("success"))
	if err != nil {
		return &ErrorType{
			text: err.Error(), statusCode: http.StatusBadRequest,
		}
	}
	return nil
}

func GetMetricsValue(
	w http.ResponseWriter,
	r *http.Request,
	t string,
	fn func(name string) (*storage.CounterMetric, *storage.GaugeMetric, error),
) *ErrorType {
	metName, extra := GetMetricData(r)
	if extra != "" {
		return &ErrorType{
			text: "can't find parameters", statusCode: http.StatusNotFound,
		}
	}
	counter, gauge, err := fn(metName)
	if err != nil {
		return &ErrorType{
			text: "can't find parameters", statusCode: http.StatusBadRequest,
		}
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var localErr ErrorType
	switch t {
	case "counter":
		_, err = fmt.Fprint(w, counter.Value)
		localErr = ErrorType{
			text: err.Error(), statusCode: http.StatusBadRequest,
		}
	case "gauge":
		_, err = fmt.Fprint(w, gauge.Value)
		localErr = ErrorType{
			text: err.Error(), statusCode: http.StatusBadRequest,
		}
	default:
		localErr = ErrorType{
			text: "can't find parameters", statusCode: http.StatusBadRequest,
		}
	}
	return &localErr
}

// The following three functions are handlers for working with metrics.
// 1. The first handler serves a list of metrics by responding to the request with an HTML file.

func MerticsListHandler(w http.ResponseWriter, r *http.Request) {
	mem := storage.GetMemStorage()
	t, err := template.New("webpage").Parse(constants.Tmpl)

	if err != nil {
		mem.GetLogger().Panic("Metric List Handler", zap.String("error", err.Error()))
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := t.Execute(w, mem); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// The next two handlers use a map to iterate over metric methods,
// avoiding code duplication with the same logic.

// They iterate through different methods with the same arguments and return types.
// Apart from the specific method, the handlers function identically.

func MetricsUpdateHandler(r chi.Router) {
	UpdateRequests := map[string]func(string, string) error{
		"counter": server.MetricsCounterSet,
		"gauge":   server.MetricsGaugeSet,
	}
	for k := range UpdateRequests {
		r.Route(fmt.Sprintf("/%s", k), func(r chi.Router) {
			r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
				if fn, exist := UpdateRequests[k]; exist {
					if err := SetMeticsUpdate(w, r, fn); err != nil {
						http.Error(w, err.text, err.statusCode)
					}
				} else {
					http.Error(w, "Can't find parameters", http.StatusBadRequest)
				}
			})
		})
	}

	r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unknown type", http.StatusBadRequest)
	})
}

func MetricsValueHandler(r chi.Router) {
	ValueRequests := map[string]func(name string) (*storage.CounterMetric, *storage.GaugeMetric, error){
		"counter": server.MetricsCounterGet,
		"gauge":   server.MetricsGaugeGet,
	}
	for k := range ValueRequests {
		r.Route(fmt.Sprintf("/%s", k), func(r chi.Router) {
			r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				if fn, exist := ValueRequests[k]; exist {
					if err := GetMetricsValue(w, r, k, fn); err != nil {
						http.Error(w, err.text, err.statusCode)
					}
				} else {
					http.Error(w, "Can't find parameters", http.StatusBadRequest)
				}
			})
		})
	}
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unknown type", http.StatusBadRequest)
	})
}

func Webhook() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", MerticsListHandler)
		r.Route("/update", MetricsUpdateHandler)
		r.Route("/value", MetricsValueHandler)
	})

	return r
}
