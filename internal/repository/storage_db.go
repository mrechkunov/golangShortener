package repository

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

type DB struct {
	dbconn  *sql.DB
	counter int
}

func NewDB() *DB {
	db, err := NewConnect()
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
	var cnt int
	_ = db.QueryRow("select uuid from storage order by uuid desc limit 1").Scan(
		&cnt)
	cnt++
	var retDB = &DB{
		dbconn:  db,
		counter: cnt,
	}
	return retDB
}

func (d *DB) SetData(shortURL string, originalURL string) error {
	err := d.dbconn.Ping()
	if err != nil {
		logger.Log.Fatal(err)
	}
	// проверяем есть ли такой URL в DB
	var shortURLFromDB string
	d.dbconn.QueryRow("select shorturl from storage where shorturl=$1", shortURL).Scan(&shortURLFromDB)
	if shortURLFromDB == shortURL {
		logger.Log.Infoln("shortURL already exist in DB")
		return errors.New("409 Conflict")
	} else {
		sqlStatement := `INSERT INTO storage (uuid, originalurl, shorturl)
		VALUES ($1, $2, $3)`
		_, err := d.dbconn.Exec(sqlStatement, d.counter, originalURL, shortURL)
		if err != nil {
			logger.Log.Errorln("error while insert to db", err)
			return err
		}
		d.counter++
	}
	return nil
}

func (d *DB) GetData(shortURL string) (string, bool) {

	err := d.dbconn.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	var res string

	err = d.dbconn.QueryRow("select originalurl from storage where shorturl=$1", shortURL).Scan(
		&res)
	isFound := true
	if err != nil {
		isFound = false
	}
	return res, isFound
}
func (d *DB) Close() error {
	return d.dbconn.Close()
}
