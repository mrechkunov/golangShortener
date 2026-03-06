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

func JSONBatchPostHandler(w http.ResponseWriter, r *http.Request) {
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	// разобрать запрос в структуру
	var requestBatch []model.RequestDataBatch
	if err := json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// проверить что пришел батч не пустой если пустой, ответить badrequest

	// пройтись по структуре, сократить ссылки в новую структуру, добавить данные в хранилище
	var responseBatch []model.ResponseDataBatch
	for _, reqBatchElement := range requestBatch {
		var resBatchElement model.ResponseDataBatch
		resBatchElement.CorrelationID = reqBatchElement.CorrelationID
		hash := sha256.Sum256([]byte(reqBatchElement.OriginalURL))
		resBatchElement.ShortURL = baseResultAdress + "/" + hex.EncodeToString(hash[:4]) // 4 байта хеша = 8 символов в hex

		err := repository.GetStorage().SetData(hex.EncodeToString(hash[:4]), reqBatchElement.OriginalURL)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
		responseBatch = append(responseBatch, resBatchElement)
	}
	// формируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// записываем ответ
	json.NewEncoder(w).Encode(responseBatch)
}

// подготовить и отправить ответ

//сокращаем url для всего батча
// for req := range reqs{
// hash := sha256.Sum256([]byte(req.original_url))
// ShortURL := baseResultAdress + "/" + hex.EncodeToString(hash[:4]) // 4 байта хеша = 8 символов в hex

// var resp []model.ResponseData
// resp.ShortURL = ShortURL
// // пишем в хранилище
// err := repository.GetStorage().SetData(hex.EncodeToString(hash[:4]), req.URL)
// if err != nil {
// 	w.WriteHeader(http.StatusBadRequest)
// 	return
// }
