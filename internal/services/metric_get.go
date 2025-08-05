// Package services
package services

import (
	"fmt"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
)

func MetricsCounterGet(name string) (*models.CounterMetric, *models.GaugeMetric, error) {
	mem := storage.GetMemStorage()
	mem.Mu.Lock()
	defer mem.Mu.Unlock()
	fmt.Println("counter", mem.Counter, name)
	if storage.ServerFlags.DataBaseDsn != "" {
		oldMet, err := GetDBCounter(name)
		return oldMet, nil, err
	}
	oldMet, err := mem.GetCounterMetric(name)

	fmt.Println("oldMet", oldMet, err)
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
