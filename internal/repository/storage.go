package repository

import (
	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
)

type StorageI interface {
	GetData(shortURL string) (string, bool)
	SetData(shortURL string, originalURL string, cookie string) error
	IsCookieExist(cookie string) bool
	GetDataByUID(uid uint32) []model.ResponseDataBatchByCookie
	IsDeleted(shortURL string) bool
	// TODO функция проставления признака удаления по shortURL
	Close() error
}

var storage StorageI

func StorageInit() StorageI {
	if config.ConfigAdreses.DBConnStr != "" {
		storage = NewDB()
		logger.Log.Infoln("work with DB:", config.ConfigAdreses.DBConnStr)
	} else if config.ConfigAdreses.JSONFile != "" {
		storage = NewSafeMapFile()
		logger.Log.Infoln("work with file:", config.ConfigAdreses.JSONFile)
	} else {
		storage = NewSafeMap()
		logger.Log.Infoln("work with memory")
	}
	return storage
}
func GetStorage() StorageI {
	return storage
}
