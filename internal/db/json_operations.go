// Package db
package db

import (
	"database/sql"
	"fmt"
	pb "github.com/rtmelsov/metrigger/proto"
	"strconv"
)

func BeginTransaction() (*sql.Tx, error) {
	DB, err := GetDataBase()
	if err != nil {
		return nil, err
	}
	return DB.Begin()
}

func CloseStmt(stmt *sql.Stmt) error {
	if err := stmt.Close(); err != nil {
		return err
	}
	return nil
}

func RollbackTx(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}

func InsertMetric(v *pb.Metric, setGauge, setCounter *sql.Stmt) error {
	switch v.MType {
	case "gauge":
		_, err := setGauge.Exec(v.ID, v.MType, v.Value)
		return err
	case "counter":
		_, err := setCounter.Exec(v.ID, v.MType, v.Delta)
		return err
	default:
		return fmt.Errorf("unsupported metric type: %s", v.MType)
	}
}

func FetchUpdatedMetric(v *pb.Metric, getGaugeCommand, getDeltaCommand *sql.Stmt) (*pb.Metric, error) {
	if v.MType == "gauge" {
		var value string
		err := getGaugeCommand.QueryRow(v.MType, v.ID).Scan(&value)
		if err != nil {
			return nil, err
		}
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return &pb.Metric{ID: v.ID, MType: v.MType, Value: f}, nil
	}

	var delta string
	err := getDeltaCommand.QueryRow(v.MType, v.ID).Scan(&delta)
	if err != nil {
		return nil, err
	}
	d, err := strconv.ParseInt(delta, 10, 64)
	if err != nil {
		return nil, err
	}
	return &pb.Metric{ID: v.ID, MType: v.MType, Delta: d}, nil
}
