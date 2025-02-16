package services

import (
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"strconv"
)

func MetricsCounterGet(name string) (*models.CounterMetric, *models.GaugeMetric, error) {
	mem := storage.GetMemStorage()
	mem.Mu.Lock()
	defer mem.Mu.Unlock()
	if storage.ServerFlags.DataBaseDsn != "" {
		oldMet, err := GetDBCounter(name)
		return oldMet, nil, err
	}
	oldMet, err := mem.GetCounterMetric(name)
	return oldMet, nil, err
}

func MetricsGaugeGet(name string) (*models.CounterMetric, *models.GaugeMetric, error) {
	mem := storage.GetMemStorage()
	mem.Mu.Lock()
	defer mem.Mu.Unlock()
	if storage.ServerFlags.DataBaseDsn != "" {
		oldMet, err := GetDBGauge(name)
		return nil, oldMet, err
	}
	oldMet, err := mem.GetGaugeMetric(name)
	return nil, oldMet, err
}

func MetricsGaugeSet(name string, val string) error {
	mem := storage.GetMemStorage()
	mem.Mu.Lock()
	defer mem.Mu.Unlock()
	if storage.ServerFlags.DataBaseDsn != "" {
		mem.GetLogger().Info("db set info: ", zap.String("DataBaseDsn", storage.ServerFlags.DataBaseDsn), zap.String("name", name), zap.String("val", val))
		return SetDBGauge(name, val)
	}
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
	mem.Mu.Lock()
	defer mem.Mu.Unlock()
	if storage.ServerFlags.DataBaseDsn != "" {
		mem.GetLogger().Info("db set info: ", zap.String("DataBaseDsn", storage.ServerFlags.DataBaseDsn), zap.String("name", name), zap.String("val", val))
		return SetDBGounter(name, val)
	}
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
