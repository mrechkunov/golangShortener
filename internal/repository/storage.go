package repository

import (
	"github.com/mrechkunov/golangShortener.git/internal/config"
)

type Stor interface {
	GetData(shortURL string) (string, bool)
	SetData(shortURL string, originalURL string)
}

// repository.Storage.ReadDataFromFile()
// var Storage *repository.SafeMap = repository.NewSafeMap()
func StorageInit() Stor {
	var Storage Stor
	if config.ConfigAdreses.DBConnStr != "" {
		//var Storage *SafeMapDB = NewSafeMapDB()
		//Storage = NewSafeMapFile()
	} else if config.ConfigAdreses.JSONFile != "" {
		var Storage = NewSafeMapFile()
		Storage.ReadDataFromFile()
	} else {
		Storage = NewSafeMap()
	}
	return Storage
}

var Storage = StorageInit()
