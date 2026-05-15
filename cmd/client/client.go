package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func postJSON() {
	url := "http://localhost:8080/api/shorten"
	jsonData, err := json.Marshal(model.RequestData{URL: "http://yandex.ru/" + randomString(10)})
	if err != nil {
		logger.Log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log.Infoln(err)
	}
	req.Header.Set("Connection", "close") // очень важно сказать серверу закрыть за нами дврецу
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Infoln(err)
	}
	req.Body.Close()
	resp.Body.Close()
}

func main() {

	for range 100000 {
		postJSON()
	}
}
