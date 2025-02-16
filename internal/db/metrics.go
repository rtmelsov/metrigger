package db

import (
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
)

func GetMetric(key string, name string) (string, error) {
	log := storage.GetMemStorage().GetLogger()
	row := db.QueryRow(`
		SELECT metric_value from metrics where metric_type = $1 and metric_name = $2
	`, key, name)
	var value string
	err = row.Scan(&value)

	log.Info("check get method", zap.String("key", key), zap.String("name", name), zap.Error(err))
	return value, nil
}

func SetMetric(key string, name string, value string) error {
	log := storage.GetMemStorage().GetLogger()
	url := `
		INSERT INTO metrics (metric_name, metric_type, metric_value)
        VALUES ($1, $2, $3)
        ON CONFLICT (metric_name, metric_type) 
        DO UPDATE SET metric_value = EXCLUDED.metric_value;
	`
	if key == "counter" {
		url = `
			INSERT INTO metrics (metric_name, metric_type, metric_value)
        	VALUES ($1, $2, $3)
        	ON CONFLICT (metric_name, metric_type) 
        	DO UPDATE SET metric_value = metrics.metric_value + EXCLUDED.metric_value;
		`
	}
	_, err = db.Exec(url, name, key, value)
	log.Info("check set method", zap.String("key", key), zap.String("name", name), zap.Error(err))
	rows, err := db.Query("SELECT id, metric_name, metric_type, metric_value FROM metrics")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		// Читаем данные из таблицы и печатаем в консоль
		var (
			ID          int
			MetricName  string
			MetricType  string
			MetricValue string
		)
		rows.Scan(&ID, &MetricName, &MetricType, &MetricValue)
		log.Info("check set method", zap.String("MetricName", MetricName), zap.String("MetricValue", MetricValue), zap.String("MetricType", MetricType))
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	log.Info("checked set method")
	return err
}
