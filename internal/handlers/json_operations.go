package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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
	resp, err := helpers.JSONElementParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var metrics []interface{}

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
}

func JSONUpdate(w http.ResponseWriter, r *http.Request) {
	storage.GetMemStorage().GetLogger().Info("in update handler")
	resp, err := helpers.JSONElementParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metric, err, statusCode := update(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v: %v", "error while updating", err.Error()), statusCode)
		return
	}
	data, err := json.Marshal(*metric)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v: %v", "failed to Marshaling JSON", err.Error()), http.StatusInternalServerError)
		return
	}
	err = SendData(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v: %v", "failed to Marshaling JSON", err.Error()), http.StatusInternalServerError)
	}
}

func JSONUpdateList(w http.ResponseWriter, r *http.Request) {
	response, err := helpers.JSONListParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if storage.ServerFlags.DataBaseDsn != "" {
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
	}
	var dataList []interface{}
	for _, v := range *response {
		metric, err, statusCode := update(&v)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v: %v", "error while updating", err.Error()), statusCode)
			return
		}
		dataList = append(dataList, *metric)
	}
	data, err := json.Marshal(dataList)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v: %v", "failed to Marshaling JSON", err.Error()), http.StatusInternalServerError)
		return
	}
	err = SendData(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v: %v", "failed to Marshaling JSON", err.Error()), http.StatusInternalServerError)
	}
}

func update(resp *models.Metrics) (*interface{}, error, int) {
	var fn func(string, string) error

	var val string

	switch resp.MType {
	case "counter":
		storage.GetMemStorage().GetLogger().Info("first in update")
		if resp.Delta == nil {
			return nil, errors.New("error while checking delta"), http.StatusNotFound
		}
		val = strconv.Itoa(int(*resp.Delta))
		fn = services.MetricsCounterSet
	case "gauge":

		storage.GetMemStorage().GetLogger().Info("second in update")
		if resp.Value == nil {
			return nil, errors.New("error while checking value"), http.StatusNotFound
		}
		val = strconv.FormatFloat(*resp.Value, 'f', -1, 64)
		fn = services.MetricsGaugeSet
	default:
		return nil, errors.New("can't find type"), http.StatusNotFound
	}

	aliasErr := SetMeticsUpdate(resp.ID, val, fn)
	if aliasErr != nil {
		return nil, errors.New(aliasErr.Text), aliasErr.StatusCode
	}

	var metric interface{}
	if resp.MType == "counter" {
		obj, _, err := services.MetricsCounterGet(resp.ID)
		if err != nil {
			return nil, errors.New("failed to find element"), http.StatusInternalServerError
		}
		num := int64(obj.Value)
		resp.Delta = &num
		metric = resp
	} else {
		_, obj, err := services.MetricsGaugeGet(resp.ID)
		if err != nil {
			return nil, errors.New("failed to find element"), http.StatusInternalServerError
		}
		resp.Value = &obj.Value
		metric = resp
	}
	return &metric, nil, 0
}
