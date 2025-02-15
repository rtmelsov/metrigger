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
	once.Do(func() {
		m := storage.GetMemStorage()
		//db, err = sql.Open("postgres", "test:test@/dbname")
		db, err = sql.Open("postgres", storage.ServerFlags.DataBaseDsn)

		if err != nil {
			m.GetLogger().Panic("error while connecting to db", zap.String("error", err.Error()))
			return
		}

		if err = db.Ping(); err != nil {
			m.GetLogger().Panic("error while ping db", zap.String("error", err.Error()))
			return
		}
	})

	return db, err
}
