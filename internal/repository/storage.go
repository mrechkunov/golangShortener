package repository

import (
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

type storageI interface {
	GetData(shortURL string) (string, bool)
	SetData(shortURL string, originalURL string)
}

func StorageInit() storageI {
	
	config.Init()
	var Storage storageI
	if config.ConfigAdreses.DBConnStr != "" {
		Storage = NewDB()
		logger.Log.Infoln("work with DB:", config.ConfigAdreses.DBConnStr)
	} else if config.ConfigAdreses.JSONFile != "" {
		Storage = NewSafeMapFile()
		logger.Log.Infoln("work with file:", config.ConfigAdreses.JSONFile)
	} else {
		Storage = NewSafeMap()
		logger.Log.Infoln("work with memory")
	}
	return Storage
}

var Storage = StorageInit()
