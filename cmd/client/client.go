package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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

func main() {

	for i := 0; i < 100000; i++ {
		url := "http://localhost:8080/api/shorten"

		// 3. Маршалинг (преобразование) строки в JSON
		//str := randomString(10)

		str := "http://yandex.ru/" + randomString(10)

		data := model.RequestData{URL: str}
		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Log.Fatal(err)
		}
		fmt.Println(str)
		// 4. Отправляем POST запрос с телом в формате JSON
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		cookie := &http.Cookie{
			Name:  "shorterner",
			Value: "e7028aa8df600638dfdb11f3bb96eacc0ad6c8beb1bf70a2c7395f4c087224c2bb82207e",
		}
		req.AddCookie(cookie)

		// Create a transport that disables automatic compression
		tr := &http.Transport{
			DisableCompression: true,
		}

		// Assign the transport to a custom client
		client := &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}
		resp, _ := client.Do(req)

		// resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		// // 2. Read the body into a byte slice
		// bodyBytes, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	// handle error
		// }

		// // 3. Convert to string if needed
		// bodyString := string(bodyBytes)

		// Выводим статус ответа
		fmt.Println("Status:", resp.Status)
		fmt.Println(i)
		//fmt.Println("Cookie:", resp.Cookies())

		resp.Body.Close()
	}
}
