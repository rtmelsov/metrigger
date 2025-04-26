package interfaces

import (
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
)

// Getter интерфейс для работы с методами для получения данных со стороны сервера
type Getter interface {
	GetLogger() *zap.Logger                                      // Получение метода для логирования
	GetGaugeMetric(name string) (*models.GaugeMetric, error)     // Получение Gauge метрики
	GetCounterMetric(name string) (*models.CounterMetric, error) // Получение Counter метрики
}

// Setter интерфейс для работы с методами для получения данных со стороны сервера
type Setter interface {
	SetGaugeMetric(name string, value models.GaugeMetric)     // Запись Gauge метрики
	SetCounterMetric(name string, value models.CounterMetric) // Запись Counter метрики
	SetDataToFile(value models.GaugeMetric) error             // Запись данных в файл
}
