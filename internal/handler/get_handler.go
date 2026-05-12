package handler

import (
	"net/http"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

const cookieName = "shorterner"
const cookieTTL = 24 * time.Hour

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	shortURL := string(req.RequestURI)
	shortURL = shortURL[1:]
	longURL, isFound := repository.GetStorage().GetData(shortURL)
	if repository.GetStorage().IsDeleted(shortURL) {
		http.Error(res, "short URL is deleted", http.StatusGone)
		return
	}
	if !isFound {
		http.Error(res, "short URL not found", http.StatusBadRequest)
		return
	}
	// направляем на аудит

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
	uid, err := cryptoauth.GetIDFromCookie(cookie.Value)
	if err != nil {
		logger.Log.Infoln(err)
		uid = 0
	}

	event := model.ObserverEvent{
		Ts:          time.Now(),
		Action:      "follow",
		UserId:      uid,
		OriginalURL: longURL,
	}

	go config.PublisherAudit.Event(event)

	res.Header().Set("Location", longURL)
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
