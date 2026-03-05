package repository

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

type DB struct {
	db      *sql.DB
	counter int
}

func NewDB() *DB {
	db, err := NewConnect()
	if err != nil {
		logger.Log.Errorln("error while db connection")
	}
	defer db.Close()

	migrationsPath := "file://../../migrations"

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
	var cnt int
	_ = db.QueryRow("select uuid from storage order by uuid desc limit 1").Scan(
		&cnt)
	// if err != nil {
	// 	logger.Log.Infoln("Error while select last id from DB", err)
	// }
	cnt++
	var retDB = DB{
		db:      db,
		counter: cnt,
	}
	return &retDB
}

func (d *DB) SetData(shortURL string, originalURL string) {
	db, err := NewConnect()
	if err != nil {
		logger.Log.Errorln("error while db connection")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Log.Fatal(err)
	}
	logger.Log.Infoln("Successfully connected to the database!")

	sqlStatement := `INSERT INTO storage (uuid, originalurl, shorturl)
		VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, d.counter, originalURL, shortURL)
	if err != nil {
		logger.Log.Errorln("error while insert to db", err)
	}
	d.counter++

}

func (s *DB) GetData(shortURL string) (string, bool) {
	db, err := NewConnect()
	if err != nil {
		logger.Log.Errorln("error while db connection")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Log.Fatal(err)
	}
	var res string
	logger.Log.Infoln("Successfully connected to the database!")
	err = db.QueryRow("select originalurl from storage where shorturl=$1", shortURL).Scan(
		&res)
	isFound := true
	if err != nil {
		isFound = false
	}
	return res, isFound
}
