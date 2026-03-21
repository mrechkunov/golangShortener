package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
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

	m, err := migrate.New(
		config.ConfigAdreses.MigrationsPath,
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

	// set DB struct (counter uuid & db conn)
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

func (d *DB) SetData(shortURL string, originalURL string, cookie string) error {
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
		uid, _ := cryptoauth.GetIDFromCookie(cookie)
		sqlStatement := `INSERT INTO storage (count, uuid, originalurl, shorturl, cookie) VALUES ($1, $2, $3, $4, $5)`
		_, err := d.dbconn.Exec(sqlStatement, d.counter, uid, originalURL, shortURL, cookie)
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
	err = d.dbconn.QueryRow("select originalurl from storage where shorturl=$1", shortURL).Scan(&res)
	isFound := true
	if err != nil {
		isFound = false
	}
	return res, isFound
}
func (d *DB) Close() error {
	return d.dbconn.Close()
}

func (d *DB) IsCookieExist(cookie string) bool {
	err := d.dbconn.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	var isFound bool
	var res string
	err = d.dbconn.QueryRow("select * from storage where cookie=$1", cookie).Scan(&res)
	if res == cookie {
		isFound = true
	} else if err != nil {
		isFound = false
	}
	return isFound
}

func (d *DB) GetDataByUID(uid uint32) []model.ResponseDataBatchByCookie {
	var result []model.ResponseDataBatchByCookie

	rows, err := d.dbconn.Query("SELECT shortURL, originalURL FROM storage where uuid=$1", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r model.ResponseDataBatchByCookie
		if err := rows.Scan(&r.ShortURL, &r.OriginalURL); err != nil {
			log.Fatal(err)
		}
		result = append(result, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}
