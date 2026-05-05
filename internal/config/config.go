package config

import (
	"flag"
	"os"
)

type Adreses struct {
	ServerBindAdress   string
	ResultServerAdress string
	MigrationsPath     string
	JSONFile           string
	DBConnStr          string
	AuditFile          string
	AuditUrl           string
}

var ConfigAdreses = Adreses{
	ServerBindAdress:   "localhost:8080",
	ResultServerAdress: "http://localhost:8080", // для работы unit теста
	JSONFile:           "",
	DBConnStr:          "",
	AuditFile:          "",
	AuditUrl:           "",
}

func Init() {
	ba := flag.String("a", "localhost:8080", "adress to server run")
	ra := flag.String("b", "http://localhost:8080", "default responce server adress")
	mp := flag.String("m", "file://migrations", "default migration PATH")
	jf := flag.String("f", "", "default storage file")
	cs := flag.String("d", "", "default DBConnStr")
	af := flag.String("audit-file", "", "default audit file")
	au := flag.String("audit-url", "", "default audit url")
	//jf := flag.String("f", "file.txt", "default storage file")
	//cs := flag.String("d", "postgres://yapra:yaprapass@10.254.40.123:5432/yandexpracticum?sslmode=disable", "default DBConnStr")
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

	if migratoinsPath, isEnvMigrationsPath := os.LookupEnv("MIGRATIONS_PATH"); isEnvMigrationsPath {
		ConfigAdreses.MigrationsPath = migratoinsPath
	} else {
		ConfigAdreses.MigrationsPath = *mp
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
	ConfigAdreses.AuditFile = *af
	ConfigAdreses.AuditUrl = *au
}
