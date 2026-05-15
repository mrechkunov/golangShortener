package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

// JSONPostHandler shorting url from json and insert it in DB
func JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	baseResultAdress := config.ConfigAdreses.ResultServerAdress
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	//проверяем cookie если нет/не проходит проверку, выдаем новую

	var isExist, isValid bool
	var cookie *http.Cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		logger.Log.Infoln("cookie is not exist", err)
		isExist = false
	} else {
		isExist = repository.GetStorage().IsCookieExist(cookie.Value)
		isValid, err = cryptoauth.ValidateCookieSign(cookie.Value)
		if err != nil {
			logger.Log.Infoln("cookie is not valid", err)
		}
	}
	if !isExist || !isValid {
		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    cryptoauth.GenerateNewCookie(),
			Expires:  time.Now().Add(cookieTTL),
			HttpOnly: true,
		}
	}
	http.SetCookie(w, cookie)
	var req model.RequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//сокращаем url
	hash := sha256.Sum256([]byte(req.URL))
	shortstr := hex.EncodeToString(hash[:4])
	ShortURL := baseResultAdress + "/" + shortstr // 4 байта хеша = 8 символов в hex
	var resp model.ResponseData
	resp.ShortURL = ShortURL
	// пишем в хранилище
	err = repository.GetStorage().SetData(shortstr, req.URL, cookie.Value)
	if err != nil {
		//формируем заголовок ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		//записываем ответ
		json.NewEncoder(w).Encode(resp)
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
		OriginalURL: req.URL,
	}
	go config.PublisherAudit.Event(event)

	//формируем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//записываем ответ
	json.NewEncoder(w).Encode(resp)

}
