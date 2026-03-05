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
	shortURL := string(req.RequestURI)
	shortURL = shortURL[1:]
	longURL, isFound := repository.GetStorage().GetData(shortURL)
	if !isFound {
		http.Error(res, "short URL not found", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", longURL)
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
