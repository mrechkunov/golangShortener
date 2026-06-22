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
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

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
		logger.Log.Warnln("Ошибка при парсинге CIDR:", err)
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
	json.NewEncoder(res).Encode(responseStatData)
}
