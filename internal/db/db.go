package db

import (
	"database/sql"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once sync.Once
	db   *sql.DB
	err  error
)

func GetDataBase() (*sql.DB, error) {
	err = nil
	once.Do(func() {
		m := storage.GetMemStorage()
		//db, err = sql.Open("postgres", "test:test@/dbname")
		db, err = sql.Open("postgres", storage.ServerFlags.DataBaseDsn)

		if err != nil {
			m.GetLogger().Panic("error while connecting to db", zap.String("error", err.Error()))
			return
		}

		if err == nil {
			if err = db.Ping(); err != nil {
				m.GetLogger().Panic("error while ping db", zap.String("error", err.Error()))
				return
			}
		}

		if err == nil {
			_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS metrics (
				id SERIAL PRIMARY KEY,
				metric_name TEXT NOT NULL,
				metric_type TEXT NOT NULL,
				metric_value DOUBLE PRECISION NOT NULL,
				UNIQUE (metric_name, metric_type)  -- Запрещает дубликаты в этих двух колонках
			);
		`)
		}
	})

	return db, err
}
