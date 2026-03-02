package db

import (
	"database/sql"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

func Connect() (*sql.DB, error) {
	connStr := config.ConfigAdreses.DBConnStr
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Log.Errorln(err)
	}
	return db, err
}
