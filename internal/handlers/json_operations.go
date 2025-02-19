package handlers

import (
	"encoding/json"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"net/http"
	"strconv"

	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/services"
)

func SendData(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}

func JSONGet(w http.ResponseWriter, r *http.Request) {
	response, err := helpers.JSONParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var metrics []interface{}

	if len(*response) > 1 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	for _, resp := range *response {
		var fn func(name string) (*models.CounterMetric, *models.GaugeMetric, error)
		switch resp.MType {
		case "counter":
			fn = services.MetricsCounterGet
		case "gauge":
			fn = services.MetricsGaugeGet
		default:
			storage.GetMemStorage().GetLogger().Info("first in resp.MType check in get")
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
		metrics = append(metrics, metric)
	}
	var data []byte
	if len(metrics) == 1 {
		data, err = json.Marshal(metrics[0])
	} else {
		data, err = json.Marshal(metrics)
	}
	if err != nil {
		http.Error(w, "Failed to Marshal JSON", http.StatusInternalServerError)
		return
	}
	err = SendData(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func JSONUpdate(w http.ResponseWriter, r *http.Request) {
	storage.GetMemStorage().GetLogger().Info("in update handler")
	response, err := helpers.JSONParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var metrics []interface{}

	if storage.ServerFlags.DataBaseDsn != "" && len(*response) > 1 {
		storage.GetMemStorage().GetLogger().Info("in db")
		updatedMetrics, err := UpdateMetrics(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(*updatedMetrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = SendData(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	} else {
		storage.GetMemStorage().GetLogger().Info("not in db")
	}

	for _, resp := range *response {
		var fn func(string, string) error

		var val string

		switch resp.MType {
		case "counter":
			storage.GetMemStorage().GetLogger().Info("first in update")
			if resp.Delta == nil {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			val = strconv.Itoa(int(*resp.Delta))
			fn = services.MetricsCounterSet
		case "gauge":

			storage.GetMemStorage().GetLogger().Info("second in update")
			if resp.Value == nil {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			val = strconv.FormatFloat(*resp.Value, 'f', -1, 64)
			fn = services.MetricsGaugeSet
		default:

			storage.GetMemStorage().GetLogger().Info("default in update")
			http.Error(w, "", http.StatusNotFound)
			return
		}

		aliasErr := SetMeticsUpdate(resp.ID, val, fn)
		if aliasErr != nil {
			http.Error(w, aliasErr.Text, aliasErr.StatusCode)
			return
		}

		var metric interface{}
		if resp.MType == "counter" {
			obj, _, err := services.MetricsCounterGet(resp.ID)
			if err != nil {
				http.Error(w, "Failed to find element", http.StatusInternalServerError)
				return
			}
			num := int64(obj.Value)
			resp.Delta = &num
			metric = resp
		} else {
			_, obj, err := services.MetricsGaugeGet(resp.ID)
			if err != nil {
				http.Error(w, "Failed to find element", http.StatusInternalServerError)
				return
			}
			resp.Value = &obj.Value
			metric = resp
		}
		metrics = append(metrics, metric)

	}
	var data []byte
	if len(metrics) == 1 {
		data, err = json.Marshal(metrics[0])
	} else {
		data, err = json.Marshal(metrics)
	}
	if err != nil {
		http.Error(w, "Failed to Marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
