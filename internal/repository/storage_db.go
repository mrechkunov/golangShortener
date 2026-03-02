package repository

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/config/db"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

func SetData(shortURL string, originalURL string) {
	db, err := db.NewConnect()
	//_, err := db.NewConnect()
	if err != nil {
		logger.Log.Errorln("error while db connection")
	}
	migrationsPath := "file://migrations"

	m, err := migrate.New(
		migrationsPath,
		config.ConfigAdreses.DBConnStr,
	)
	if err != nil {
		logger.Log.Fatalln("Error initializing migrate:", err)
	}
	// Apply all available migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatalln("Error applying migrations:", err)
	}

	logger.Log.Infoln("Database migrations applied successfully!")
	db.Ping()
}

// func GetData(shortURL string) (string, bool) {

// }
