package config

import (
	"flag"
	"os"
)

type Adreses struct {
	ServerBindAdress   string
	ResultServerAdress string
}

var ConfigAdreses = Adreses{
	ServerBindAdress:   "localhost:8080",
	ResultServerAdress: "http://localhost:8080", // для работы unit теста
}

func Init() {

	ba := flag.String("a", "localhost:8080", "adress to server run")
	ra := flag.String("b", "http://localhost:8080", "default responce server adress")
	flag.Parse()

	var isEnvBindSrv bool
	var isEnvResSrv bool
	// проверяем установленны ли переменные окружения
	_, isEnvBindSrv = os.LookupEnv("SERVER_ADDRESS")
	_, isEnvResSrv = os.LookupEnv("BASE_URL")
	// если переиенные окружения установленны, берем адреса из них, иначе берем адреса из флага
	if isEnvBindSrv {
		ConfigAdreses.ServerBindAdress = os.Getenv("SERVER_ADDRESS")
	} else {
		ConfigAdreses.ServerBindAdress = *ba
	}
	if isEnvResSrv {
		ConfigAdreses.ResultServerAdress = os.Getenv("BASE_URL")
	} else {
		ConfigAdreses.ResultServerAdress = *ra
	}
}
