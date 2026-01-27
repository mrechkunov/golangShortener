package handler

import (
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	shortUrl := string(req.RequestURI)
	shortUrl = shortUrl[1:]
	res.Header().Set("Location", repository.SelectData(shortUrl))
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
