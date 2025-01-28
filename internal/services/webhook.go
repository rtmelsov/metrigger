package services

import (
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"strconv"
)

func MetricsCounterGet(name string) (*models.CounterMetric, *models.GaugeMetric, error) {
	mem := storage.GetMemStorage()
	oldMet, err := mem.GetCounterMetric(name)
	return oldMet, nil, err
}

func MetricsGaugeGet(name string) (*models.CounterMetric, *models.GaugeMetric, error) {
	mem := storage.GetMemStorage()
	oldMet, err := mem.GetGaugeMetric(name)
	return nil, oldMet, err
}

func MetricsGaugeSet(name string, val string) error {
	mem := storage.GetMemStorage()
	met := storage.NewGaugeMetric()
	met.Type = "gauge"
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	met.Value = f
	mem.SetGaugeMetric(name, *met)
	return nil
}

func MetricsCounterSet(name string, val string) error {
	mem := storage.GetMemStorage()
	met := storage.NewCounterMetric()
	met.Type = "counter"
	i, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	oldCount := 0
	oldMet, err := mem.GetCounterMetric(name)
	if err != nil {
		logger := storage.GetMemStorage().GetLogger()
		logger.Debug("Error: %v\r\n", zap.String("error", err.Error()))
	} else {
		oldCount = oldMet.Value
	}
	met.Value = oldCount + i
	mem.SetCounterMetric(name, *met)
	return nil
}
