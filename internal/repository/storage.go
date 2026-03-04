package repository

import (
	"fmt"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

type Stor interface {
	GetData(shortURL string) (string, bool)
	SetData(shortURL string, originalURL string)
}

// repository.Storage.ReadDataFromFile()
// var Storage *repository.SafeMap = repository.NewSafeMap()
func StorageInit() Stor {
	config.Init()
	var Storage Stor
	if config.ConfigAdreses.DBConnStr != "" {
		logger.Log.Infow("work with DB")
		Storage = NewSafeMapDB()
		//Storage = NewSafeMapFile()
	} else if config.ConfigAdreses.JSONFile != "" {
		Storage = NewSafeMapFile()
		fmt.Println("work with file", config.ConfigAdreses.JSONFile)
	} else {
		Storage = NewSafeMap()
		fmt.Println("work with memory")
	}
	return Storage
}

var Storage = StorageInit()
