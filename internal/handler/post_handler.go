package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

// PostHandler shorting url from POST request and insert it in DB
func PostHandler(res http.ResponseWriter, req *http.Request) {
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
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
			Expires:  time.Now().Add(cookieTTL),
			HttpOnly: true,
		}
	}
	http.SetCookie(res, cookie)

	//читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body reading error", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	//сокращаем url
	hash := sha256.Sum256([]byte(body))
	shortURL := baseResultAdress + "/" + hex.EncodeToString(hash[:4]) // 4 байта хеша = 8 символов в hex
	err = repository.GetStorage().SetData(hex.EncodeToString(hash[:4]), string(body), cookie.Value)
	if err != nil {
		res.Header().Set("content-type", "text/plain; charset=utf-8")
		res.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))
		res.WriteHeader(http.StatusConflict)
		//записываем ответ
		res.Write([]byte(shortURL))
		return
	}
	// направляем на аудит
	uid, err := cryptoauth.GetIDFromCookie(cookie.Value)
	if err != nil {
		logger.Log.Warnln(err)
	}
	event := model.ObserverEvent{
		Ts:          time.Now(),
		Action:      "shorten",
		UserId:      uid,
		OriginalURL: string(body),
	}
	go config.PublisherAudit.Event(event)
	//формируем заголовок ответа

	res.Header().Set("content-type", "text/plain; charset=utf-8")
	res.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))
	res.WriteHeader(http.StatusCreated)
	//записываем ответ
	res.Write([]byte(shortURL))

}
