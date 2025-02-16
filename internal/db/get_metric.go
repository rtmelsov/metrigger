package db

func GetMetric(key string, name string) (string, error) {
	row := db.QueryRow(`
		SELECT metric_value from metrics where metric_type = $1 and metric_name = $2
	`, key, name)
	var value string
	err = row.Scan(&value)
	return value, nil
}

func SetMetric(key string, name string, value string) error {
	_, err = db.Exec(`
		INSERT INTO metrics (metric_name, metric_type, metric_value) 
VALUES ($1, $2, $3);
	`, name, key, value)
	return err
}
