package handler

import (
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

func GetHandlerPingDB(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	db, err := repository.NewConnect()
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}
	db.Close()
}
