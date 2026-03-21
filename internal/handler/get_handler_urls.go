package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
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
	fmt.Println("incomign cookie:", cookie)
	if err != nil {
		logger.Log.Infoln("cookie is not exist", err)
		isExist = false
	} else {
		isExist = repository.GetStorage().IsCookieExist(cookie.Value)
		isValid, _ = cryptoauth.ValidateCookieSign(cookie.Value)
	}
	if !isValid {
		http.Error(res, "cookie is not valid", http.StatusNoContent)
		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    cryptoauth.GenerateNewCookie(),
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		return
	}
	if !isExist {
		http.Error(res, "cookie is not exist in storage", http.StatusOK)
		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    cryptoauth.GenerateNewCookie(),
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		return
	}
	uid, err := cryptoauth.GetIDFromCookie(cookie.Value)
	if err != nil {
		http.Error(res, "no ID in cookie", http.StatusUnauthorized)
		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    cryptoauth.GenerateNewCookie(),
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		return
	}

	// Выбрать из хранилища все записи с uid
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	responseBatch := repository.GetStorage().GetDataByUID(uid)
	// fmt.Println("len of responseBatch", len(responseBatch))
	// if len(responseBatch) == 0 {
	// 	fmt.Println("set 204")
	// 	res.WriteHeader(http.StatusNoContent)
	// 	fmt.Println("set coockie")
	// 	http.SetCookie(res, cookie)
	// 	return
	// }
	var result []model.ResponseDataBatchByCookie
	for _, rb := range responseBatch {
		rb.ShortURL = baseResultAdress + "/" + rb.ShortURL
		result = append(result, rb)
	}

	cookie = &http.Cookie{
		Name:     cookieName,
		Value:    cookie.Value,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(res, cookie)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	// записываем ответ
	json.NewEncoder(res).Encode(result)
}
