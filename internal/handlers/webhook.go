package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rtmelsov/metrigger/internal/server"
)

var Tmpl = `
<!DOCTYPE html>
			<html>
			<head>
				<title>Metrics</title>
			</head>
			<body>
				<div>
					<h3>Counter</h3>
					{{range $category, $product := .Counter}}
						<div>
							<h5>{{$category}}</h5>
							<ul>{{$product.Value}}</ul>
						</div>
					{{end}}
				</div>
				<div>
					<h3>Gauge</h3>
					{{range $category, $product := .Gauge}}
						<div>
							<h5>{{$category}}</h5>
							<ul>{{$product.Value}}</ul>
						</div>
					{{end}}
				</div>

			</body>
	</html>
	`

func GetMetricData(r *http.Request) (string, string) {
	paths := strings.Split(r.URL.String(), "/")
	fmt.Printf("paths: %v\r\n", paths)
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
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			mem := server.MetricsGet()
			t, err := template.New("webpage").Parse(Tmpl)
			if err != nil {
				log.Panic(err)
			}
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			if err := t.Execute(w, mem); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
		})
		r.Route("/update", func(r chi.Router) {
			r.Route("/counter", func(r chi.Router) {
				r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
					metName, metVal := GetMetricData(r)
					if metName == "" || metVal == "" {
						http.Error(w, "Can't find parameters", http.StatusNotFound)
						return
					}
					err := server.MetricsCounterSet(metName, metVal)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
					w.WriteHeader(http.StatusOK)
				})
			})
			r.Route("/gauge", func(r chi.Router) {
				r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
					metName, metVal := GetMetricData(r)
					if metName == "" || metVal == "" {
						http.Error(w, "Can't find parameters", http.StatusNotFound)
						return
					}
					err := server.MetricsGaugeSet(metName, metVal)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
					w.WriteHeader(http.StatusOK)
				})
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Route("/counter", func(r chi.Router) {
				r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
					metName, extra := GetMetricData(r)
					if metName == "" {
						http.Error(w, "Can't find parameters", http.StatusNotFound)
						return
					}
					if extra != "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					val, err := server.MetricsCounterGet(metName)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					if _, err = fmt.Fprint(w, val.Value); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
				})
			})
			r.Route("/gauge", func(r chi.Router) {
				r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
					metName, extra := GetMetricData(r)
					if metName == "" {
						http.Error(w, "Can't find parameters", http.StatusNotFound)
						return
					}
					if extra != "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					val, err := server.MetricsGaugeGet(metName)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					if _, err = fmt.Fprint(w, val.Value); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
				})
			})
		})
	})

	return r
}
