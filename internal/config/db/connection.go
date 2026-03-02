package db

import (
	"database/sql"

	_ "github.com/jackc/pgx"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

func Connect() (*sql.DB, error) {
	connStr := config.ConfigAdreses.DBConnStr
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Log.Errorln(err)
	}
	return db, err
}
