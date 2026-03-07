package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
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
	if len(requestBatch) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Infoln("request batch is empty, lenght:", len(requestBatch))
		return
	}

	// пройтись по структуре, сократить ссылки в новую структуру, добавить данные в хранилище
	var responseBatch []model.ResponseDataBatch
	var errconfl error = nil
	for _, reqBatchElement := range requestBatch {
		var resBatchElement model.ResponseDataBatch
		resBatchElement.CorrelationID = reqBatchElement.CorrelationID
		hash := sha256.Sum256([]byte(reqBatchElement.OriginalURL))
		resBatchElement.ShortURL = baseResultAdress + "/" + hex.EncodeToString(hash[:4]) // 4 байта хеша = 8 символов в hex
		if errconfl == nil {
			errconfl = repository.GetStorage().SetData(hex.EncodeToString(hash[:4]), reqBatchElement.OriginalURL)
		} else {
			repository.GetStorage().SetData(hex.EncodeToString(hash[:4]), reqBatchElement.OriginalURL)
		}
		responseBatch = append(responseBatch, resBatchElement)
	}

	if errconfl != nil {
		// формируем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		// записываем ответ
		json.NewEncoder(w).Encode(responseBatch)
		return
	}
	// формируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// записываем ответ
	json.NewEncoder(w).Encode(responseBatch)
}
