package config

import (
	"flag"
	"os"
)

type Adreses struct {
	ServerBindAdress   string
	ResultServerAdress string
	JSONFile           string
}

var ConfigAdreses = Adreses{
	ServerBindAdress:   "localhost:8080",
	ResultServerAdress: "http://localhost:8080", // для работы unit теста
	JSONFile:           "./file.txt",
}

func Init() {

	ba := flag.String("a", "localhost:8080", "adress to server run")
	ra := flag.String("b", "http://localhost:8080", "default responce server adress")
	jf := flag.String("f", "./file.txt", "default storage file")
	flag.Parse()

	var isEnvBindSrv bool
	var isEnvResSrv bool
	var isEnvJSONFile bool
	// проверяем установленны ли переменные окружения
	_, isEnvBindSrv = os.LookupEnv("SERVER_ADDRESS")
	_, isEnvResSrv = os.LookupEnv("BASE_URL")
	_, isEnvJSONFile = os.LookupEnv("FILE_STORAGE_PATH")
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
	if isEnvJSONFile {
		ConfigAdreses.JSONFile = os.Getenv("FILE_STORAGE_PATH")
	} else {
		ConfigAdreses.JSONFile = *jf
	}
}
