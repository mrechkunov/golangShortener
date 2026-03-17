package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

func NewConnect() (*sql.DB, error) {
	db, err := sql.Open("pgx", config.ConfigAdreses.DBConnStr)
	if err != nil {
		logger.Log.Errorln(err)
	}
	return db, err
}
