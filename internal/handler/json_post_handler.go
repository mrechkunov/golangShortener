package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/model"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

func JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	var req model.RequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//сокращаем url
	hash := sha256.Sum256([]byte(req.URL))
	ShortURL := baseResultAdress + "/" + hex.EncodeToString(hash[:4]) // 4 байта хеша = 8 символов в hex

	var resp model.ResponseData
	resp.ShortURL = ShortURL

	//формируем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//записываем ответ
	json.NewEncoder(w).Encode(resp.ShortURL)
	// пишем в хранилище
	repository.Storage.SetData(hex.EncodeToString(hash[:4]), string(req.URL))
}
