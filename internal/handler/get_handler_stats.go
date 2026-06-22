package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

// GetHandlerURLs return to user all urls where user is creator
func GetHandlerStats(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	//проверяем cookie если нет/не проходит проверку, выдаем новую
	var isExist, isValid bool
	var cookie *http.Cookie
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		logger.Log.Infoln("cookie is not exist in request", err)
		isExist = false
	} else {
		isExist = repository.GetStorage().IsCookieExist(cookie.Value)
		isValid, _ = cryptoauth.ValidateCookieSign(cookie.Value)
	}
	if !isValid {
		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    cryptoauth.GenerateNewCookie(),
			Expires:  time.Now().Add(cookieTTL),
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		http.Error(res, "cookie is not valid", http.StatusNoContent)
		return
	}
	if !isExist {
		logger.Log.Infoln("cookie is not exist in storage")
		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    cryptoauth.GenerateNewCookie(),
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		res.Header().Set("Content-Type", "application/json")
		http.Error(res, "cookie is not exist in storage", http.StatusOK)
		return
	}
	// uid, err := cryptoauth.GetIDFromCookie(cookie.Value)
	// if err != nil {
	// 	logger.Log.Infoln("no ID in cookie")
	// 	http.Error(res, "no ID in cookie", http.StatusUnauthorized)
	// 	cookie = &http.Cookie{
	// 		Name:     cookieName,
	// 		Value:    cryptoauth.GenerateNewCookie(),
	// 		Expires:  time.Now().Add(24 * time.Hour),
	// 		HttpOnly: true,
	// 	}
	// 	http.SetCookie(res, cookie)
	// 	res.Header().Set("Content-Type", "application/json")
	// 	return
	// }

	// Выбрать из хранилища данные статистики
	responseStatData := repository.GetStorage().GetStatData()

	cookie = &http.Cookie{
		Name:     cookieName,
		Value:    cookie.Value,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	// формируем и записываем ответ сервера
	http.SetCookie(res, cookie)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(responseStatData)
}
