package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

func DeleteHandler(c chan []string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			http.Error(res, "Only DELETE requests are allowed!", http.StatusBadRequest)
			return
		}
		//проверяем cookie если нет/не проходит проверку, выдаем новую
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
			cookie = &http.Cookie{
				Name:     cookieName,
				Value:    cryptoauth.GenerateNewCookie(),
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			}
		}
		http.SetCookie(res, cookie)

		var reqdata []string
		if err := json.NewDecoder(req.Body).Decode(&reqdata); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		// проверяем может ли пользователь помечать на удаление входящие данные
		// формируем слайс и кидаем его в канал для удаления
		var sliceToDelete []string
		for _, str := range reqdata {
			if repository.GetStorage().IsCreator(str, cookie.Value) {
				sliceToDelete = append(sliceToDelete, str)
			}
		}
		c <- sliceToDelete
		//формируем заголовок ответа
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusAccepted)
	}
}
