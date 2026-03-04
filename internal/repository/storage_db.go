package repository

import (
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
)

type SafeMapDB struct {
	mu      sync.RWMutex
	m       map[string]model.Event
	counter int
}

func NewSafeMapDB() *SafeMapDB {
	return &SafeMapDB{
		m:       make(map[string]model.Event),
		counter: 1,
	}
}

func (s *SafeMapDB) SetData(shortURL string, originalURL string) {
	db, err := NewConnect()
	//_, err := db.NewConnect()
	if err != nil {
		logger.Log.Errorln("error while db connection")
	}
	migrationsPath := "file://../../migrations"

	// // Чтение содержимого директории
	// files, err := os.ReadDir(migrationsPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Содержимое каталога:", migrationsPath)
	// for _, file := range files {
	// 	// file.IsDir() проверяет, является ли элемент каталогом
	// 	if file.IsDir() {
	// 		fmt.Println("[DIR] ", file.Name())
	// 	} else {
	// 		fmt.Println("[FILE]", file.Name())
	// 	}
	// }
	//host=10.254.40.123 user=yapra password=yaprapass dbname=yandexpracticum sslmode=disable

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

func (s *SafeMapDB) GetData(shortURL string) (string, bool) {
	return "test", true
}
