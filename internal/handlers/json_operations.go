package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/services"
	"github.com/rtmelsov/metrigger/internal/storage"
)

func jsonParse(r *http.Request) (*models.Metrics, error) {
	if r.Body == nil {
		return nil, errors.New("body is empty")
	}
	var resp *models.Metrics
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&resp)

	if err != nil {
		return nil, err
	}
	if resp.ID == "" {
		return nil, errors.New("id is empty")
	}

	return resp, nil
}

func JSONGet(w http.ResponseWriter, r *http.Request) {
	resp, err := jsonParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var fn func(name string) (*storage.CounterMetric, *storage.GaugeMetric, error)
	switch resp.MType {
	case "counter":
		fn = services.MetricsCounterGet
	case "gauge":
		fn = services.MetricsGaugeGet
	default:
		http.Error(w, "", http.StatusNotFound)
		return
	}

	var aliasErr *models.ErrorType
	counter, gauge, aliasErr := GetMetricsValue(resp.ID, "", fn)
	if aliasErr != nil {
		http.Error(w, aliasErr.Text, aliasErr.StatusCode)
		return
	}

	var metric interface{}
	if resp.MType == "counter" {
		num := int64(counter.Value)
		resp.Delta = &num
		metric = resp
	} else {
		resp.Value = &gauge.Value
		metric = resp
	}
	data, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, "Failed to Marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func JSONUpdate(w http.ResponseWriter, r *http.Request) {
	resp, err := jsonParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fn func(string, string) error

	var val string

	switch resp.MType {
	case "counter":
		if resp.Delta == nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		val = strconv.Itoa(int(*resp.Delta))
		fn = services.MetricsCounterSet
	case "gauge":
		if resp.Value == nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		val = strconv.FormatFloat(*resp.Value, 'f', -1, 64)
		fn = services.MetricsGaugeSet
	default:
		http.Error(w, "", http.StatusNotFound)
		return
	}

	aliasErr := SetMeticsUpdate(resp.ID, val, fn)
	if aliasErr != nil {
		http.Error(w, aliasErr.Text, aliasErr.StatusCode)
		return
	}

	var metric interface{}
	mem := storage.GetMemStorage()
	if resp.MType == "counter" {
		obj, err := mem.GetCounterMetric(resp.ID)
		if err != nil {
			http.Error(w, "Failed to find element", http.StatusInternalServerError)
			return
		}
		num := int64(obj.Value)
		resp.Delta = &num
		metric = resp
	} else {
		obj, err := mem.GetGaugeMetric(resp.ID)
		if err != nil {
			http.Error(w, "Failed to find element", http.StatusInternalServerError)
			return
		}
		resp.Value = &obj.Value
		metric = resp
	}

	data, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, "Failed to Marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
