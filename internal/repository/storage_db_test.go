package repository_test

import (
	"fmt"

	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

// SetData Метод для записи данных в БД
func ExampleDB_SetData() {
	// создаем сторэдж как БД
	storage := repository.NewDB()
	//вызываем метод с передачей параметров, читаем и обрабатываем ошибку
	err := storage.SetData("yandex.ru/sdfgw4sds", "yandex.ru/example", "a5e429ca94af5232196d131f95f8e4c85958500c9ecd9a140bce09680bfa94585157b851")
	if err != nil {
		logger.Log.Errorln(err)
	}
}

// GetData для чтения данных из БД по короткому URL
func ExampleDB_GetData() {
	// создаем сторэдж как БД
	storage := repository.NewDB()
	//вызываем метод с передачей параметров, читаем и обрабатываем ответ
	originalURL, isFound := storage.GetData("yandex.ru/sdfgw4sds")
	if isFound != true {
		logger.Log.Infoln("shortURL not found in DB")
	}
	// используем originalURL по назначению
	fmt.Println(originalURL)
}
