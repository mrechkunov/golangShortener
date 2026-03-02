package handler

import (
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/config/db"
)

func GetHandlerPingDB(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	_, err := db.Connect()
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}
