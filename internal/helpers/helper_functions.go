package helpers

import (
	"github.com/rtmelsov/metrigger/internal/models"
	"os"
)

func EmptyLocalStorage(path string) (*models.LocalStorage, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &models.LocalStorage{
		Gauge:   make(map[string]models.GaugeMetric),
		Counter: make(map[string]models.CounterMetric),
		File:    file,
	}, nil
}
