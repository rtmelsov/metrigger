// Package agent
package agent

import "github.com/rtmelsov/metrigger/internal/models"

// RequestToServer функция для формирования объекта для отправки в сервер
func RequestToServer(t string, key string, value float64, counter int64) *models.Metrics {
	var metric *models.Metrics

	if t == "counter" {
		metric = &models.Metrics{
			MType: t,
			ID:    key,
			Delta: &counter,
		}
	} else {
		metric = &models.Metrics{
			MType: t,
			ID:    key,
			Value: &value,
		}
	}

	return metric
}
