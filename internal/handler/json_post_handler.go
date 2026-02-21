package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/model"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

func JSONPostHandler(res http.ResponseWriter, req *http.Request) {
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	var reqq model.JSONStruct

	// Декодируем JSON из тела запроса
	err := json.NewDecoder(req.Body).Decode(&reqq)
	if err != nil {
		http.Error(res, "JSON parsing error", http.StatusBadRequest)
		return
	}

	//сокращаем url
	hash := sha256.Sum256([]byte(reqq.URL))
	reqq.ShortURL = baseResultAdress + "/" + hex.EncodeToString(hash[:4]) // 4 байта хеша = 8 символов в hex
	//формируем заголовок ответа
	res.Header().Set("content-type", "application/json")
	res.Header().Set("Content-Length", strconv.Itoa(len(reqq.ShortURL)))
	res.WriteHeader(http.StatusCreated)
	//записываем ответ
	res.Write([]byte(reqq.ShortURL))
	repository.Storage.SetData(hex.EncodeToString(hash[:4]), string(reqq.URL))
}
