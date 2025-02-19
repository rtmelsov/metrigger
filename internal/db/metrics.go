package db

import (
	"github.com/rtmelsov/metrigger/internal/constants"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
)

func GetMetric(key string, name string) (string, error) {
	logger := storage.GetMemStorage().GetLogger()
	row := db.QueryRow(constants.GetRowCommand, key, name)
	var value string
	err = row.Scan(&value)
	if err != nil {
		logger.Error("get method", zap.String("key", key), zap.String("name", name), zap.Error(err))
		return "", err
	}
	return value, nil
}

func SetMetric(key string, name string, value string) error {
	log := storage.GetMemStorage().GetLogger()
	url := constants.GaugeCommand
	if key == "counter" {
		url = constants.CounterCommand
	}
	_, err = db.Exec(url, name, key, value)
	log.Info("update method", zap.String("key", key), zap.String("name", name), zap.Error(err))

	return err
}
