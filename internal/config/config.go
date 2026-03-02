package config

import (
	"flag"
	"os"
)

type Adreses struct {
	ServerBindAdress   string
	ResultServerAdress string
	JSONFile           string
	DBConnStr          string
}

var ConfigAdreses = Adreses{
	ServerBindAdress:   "localhost:8080",
	ResultServerAdress: "http://localhost:8080", // для работы unit теста
	JSONFile:           "file.txt",
	DBConnStr:          "host=10.254.40.123 user=yapra password=yaprapass dbname=yandexpracticum sslmode=disable",
}

func Init() {
	ba := flag.String("a", "localhost:8080", "adress to server run")
	ra := flag.String("b", "http://localhost:8080", "default responce server adress")
	jf := flag.String("f", "file.txt", "default storage file")
	cs := flag.String("d", "host=10.254.40.123 user=yaPra password=YaPra dbname=yandexpracticum sslmode=disable", "default DBConnStr")
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
	if dbConnStr, isEnvDBConnStr := os.LookupEnv("DATABASE_DSN"); isEnvDBConnStr {
		ConfigAdreses.DBConnStr = dbConnStr
	} else {
		ConfigAdreses.DBConnStr = *cs
	}
}
