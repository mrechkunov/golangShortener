package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/mrechkunov/golangShortener.git/internal/config"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/repository"
)

// GetHandlerURLs return to user all urls where user is creator
func GetHandlerStats(res http.ResponseWriter, req *http.Request) {
	//если trusted subnet не задана, запрещаем вызов хэндлера
	if config.ConfigAdreses.TrustedSubnet == "" {
		http.Error(res, "No trusted subnet set", http.StatusForbidden)
		return
	}

	// проверить есть ли IP адрес в хедере X-Real-IP
	// читаем хедер
	realIP := req.Header.Get("X-Real-IP")
	// Парсим CIDR
	_, ipNet, err := net.ParseCIDR(config.ConfigAdreses.TrustedSubnet)
	if err != nil {
		logger.Log.Warnln("Error while parsing CIDR:", err)
		http.Error(res, "Error while parsing CIDR", http.StatusInternalServerError)
		return
	}
	if !ipNet.Contains(net.ParseIP(realIP)) {
		http.Error(res, "IP Address not in trusted subnet", http.StatusForbidden)
		return
	}
	// Выбрать из хранилища данные статистики
	responseStatData := repository.GetStorage().GetStatData()
	// формируем и записываем ответ сервера
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(responseStatData)
	if err != nil {
		logger.Log.Warnln("Error while encoding JSON", err)
		http.Error(res, "Error while encoding JSON", http.StatusInternalServerError)
		return
	}
}
