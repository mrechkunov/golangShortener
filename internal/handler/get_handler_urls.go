package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

func GetHandlerURLs(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	//проверяем cookie если нет/не проходит проверку, выдаем новую
	cookieName := "shorterner"
	var isExist, isValid bool
	var cookie *http.Cookie
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		logger.Log.Infoln("cookie is not exist", err)
		isExist = false
	} else {
		isExist = repository.GetStorage().IsCookieExist(cookie.Value)
		isValid, _ = cryptoauth.ValidateCookieSign(cookie.Value)
	}
	if !isExist || !isValid {
		http.Error(res, "cookie is not exist in storage", http.StatusUnauthorized)
		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    cryptoauth.GenerateNewCookie(),
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}

	}
	http.SetCookie(res, cookie)
	uid, _ := cryptoauth.GetIDFromCookie(cookie.Value)

	// Выбрать из хранилища все записи с uid
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	responseBatch := repository.GetStorage().GetDataByUID(uid)
	if len(responseBatch) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	for _, rb := range responseBatch {
		rb.ShortURL = baseResultAdress + "/" + rb.ShortURL
	}

	// формируем ответ
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	// записываем ответ
	json.NewEncoder(res).Encode(responseBatch)
}
