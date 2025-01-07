package handlers

import (
	"github.com/rtmelsov/metrigger/internal/server"
	"net/http"
	"strings"
)

func Webhook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	paths := strings.Split(r.URL.String(), "/")
	if len(paths) == 5 && paths[1] == "update" {
		metType := paths[2]
		metName := paths[3]
		metVal := paths[4]
		switch metType {
		case "counter":
			err := server.MetricsCounterSet(metName, metVal)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case "gauge":
			err := server.MetricsGaugeSet(metName, metVal)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
