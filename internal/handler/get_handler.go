package handler

import (
	"net/http"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
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
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		logger.Log.Infoln(err)
		cookie.Value = ""
	}
	uid, err := cryptoauth.GetIDFromCookie(cookie.Value)
	if err != nil {
		logger.Log.Warnln(err)
	}
	go logger.Audit("follow", uid, longURL)

	res.Header().Set("Location", longURL)
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
