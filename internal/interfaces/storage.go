package interfaces

import (
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
)

type Getter interface {
	GetLogger() *zap.Logger
	GetGaugeMetric(name string) (*models.GaugeMetric, error)
	GetCounterMetric(name string) (*models.CounterMetric, error)
}

type Setter interface {
	SetGaugeMetric(name string, value models.GaugeMetric)
	SetCounterMetric(name string, value models.CounterMetric)
	SetDataToFile(value models.GaugeMetric) error
}
