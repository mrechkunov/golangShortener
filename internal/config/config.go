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
	jf := flag.String("f", "file.txt", "default storage file")
	flag.Parse()

	// если переиенные окружения установленны, берем их, иначе берем флаг
	if serverAddress, isEnvBindSrv := os.LookupEnv("SERVER_ADDRESS"); isEnvBindSrv {
		ConfigAdreses.ServerBindAdress = serverAddress
	} else {
		ConfigAdreses.ServerBindAdress = *ba
	}
	if baseURL, isEnvResSrv := os.LookupEnv("BASE_URL"); isEnvResSrv {
		ConfigAdreses.ResultServerAdress = baseURL
	} else {
		ConfigAdreses.ResultServerAdress = *ra
	}
	if fileStoragePath, isEnvJSONFile := os.LookupEnv("FILE_STORAGE_PATH"); isEnvJSONFile {
		ConfigAdreses.JSONFile = fileStoragePath
	} else {
		ConfigAdreses.JSONFile = *jf
	}
}
