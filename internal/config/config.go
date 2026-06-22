package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

type Adreses struct {
	ServerBindAdress   string `json:"server_address"`
	ResultServerAdress string `json:"base_url"`
	MigrationsPath     string
	JSONFile           string `json:"file_storage_path"`
	DBConnStr          string `json:"database_dsn"`
	AuditFile          string
	AuditUrl           string
	HttpsEnable        bool `json:"enable_https"`
}

var Fmem *os.File
var ConfigAdreses = Adreses{
	ServerBindAdress:   "localhost:8080",
	ResultServerAdress: "http://localhost:8080", // для работы unit теста
	JSONFile:           "",
	DBConnStr:          "",
	AuditFile:          "",
	AuditUrl:           "",
}

var PublisherAudit logger.Audit
var ConfigFileAddress string
var ConfigFileData Adreses

func Init() {
	cf := flag.String("c", "", "config file address")
	ba := flag.String("a", "localhost:8080", "adress to server run")
	ra := flag.String("b", "http://localhost:8080", "default responce server adress")
	mp := flag.String("m", "file://migrations", "default migration PATH")
	jf := flag.String("f", "", "default storage file")
	cs := flag.String("d", "", "default DBConnStr")
	af := flag.String("audit-file", "", "default audit file")
	au := flag.String("audit-url", "", "default audit url")
	se := flag.Bool("s", false, "https enable")
	//jf := flag.String("f", "file.txt", "default storage file")
	//cs := flag.String("d", "postgres://yapra:yaprapass@10.254.40.123:5432/yandexpracticum?sslmode=disable", "default DBConnStr")
	flag.Parse()

	// если переиенные окружения установленны, берем их, иначе берем флаг

	if cfAddress, isConfigFileAddress := os.LookupEnv("CONFIG"); isConfigFileAddress {
		ConfigFileAddress = cfAddress
	} else {
		ConfigFileAddress = *cf
	}
	// если есть адрес конфига, парсим сначала его и присваеваем все значения структуре конфигурации
	if ConfigFileAddress != "" {
		// Считайтываем файл целиком
		data, err := os.ReadFile(ConfigFileAddress)
		if err != nil {
			logger.Log.Warnln("Ошибка чтения файла:", err)
		}
		// Распарсим JSON в структуру
		err = json.Unmarshal(data, &ConfigFileData)
		if err != nil {
			logger.Log.Warnln("Ошибка парсинга JSON:", err)
		}
	}
	if serverAddress, isEnvBindSrv := os.LookupEnv("SERVER_ADDRESS"); isEnvBindSrv {
		ConfigAdreses.ServerBindAdress = serverAddress
	} else {
		ConfigAdreses.ServerBindAdress = *ba
	}
	if ConfigAdreses.ServerBindAdress == "localhost:8080" && ConfigFileData.ServerBindAdress != "" {
		ConfigAdreses.ServerBindAdress = ConfigFileData.ServerBindAdress
	}
	if _, isEnvHttpsEnable := os.LookupEnv("ENABLE_HTTPS"); isEnvHttpsEnable {
		ConfigAdreses.HttpsEnable = true
	} else {
		ConfigAdreses.HttpsEnable = *se
	}
	if !ConfigAdreses.HttpsEnable && ConfigFileData.HttpsEnable {
		ConfigAdreses.HttpsEnable = true
	}
	if baseURL, isEnvResSrv := os.LookupEnv("BASE_URL"); isEnvResSrv {
		ConfigAdreses.ResultServerAdress = baseURL
	} else {
		ConfigAdreses.ResultServerAdress = *ra
	}
	if ConfigAdreses.ResultServerAdress == "http://localhost:8080" && ConfigFileData.ResultServerAdress != "" {
		ConfigAdreses.ResultServerAdress = ConfigFileData.ResultServerAdress
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
	if ConfigAdreses.JSONFile == "" && ConfigFileData.JSONFile != "" {
		ConfigAdreses.JSONFile = ConfigFileData.JSONFile
	}

	if dbConnStr, isEnvDBConnStr := os.LookupEnv("DATABASE_DSN"); isEnvDBConnStr {
		ConfigAdreses.DBConnStr = dbConnStr
	} else {
		ConfigAdreses.DBConnStr = *cs
	}
	if ConfigAdreses.DBConnStr == "" && ConfigFileData.DBConnStr != "" {
		ConfigAdreses.DBConnStr = ConfigFileData.DBConnStr
	}
	// создаем подписчиков
	ConfigAdreses.AuditFile = *af
	ConfigAdreses.AuditUrl = *au

	if ConfigAdreses.AuditFile != "" {
		obsFile := logger.NewObserverFile(ConfigAdreses.AuditFile)
		PublisherAudit.RegisterObserver(obsFile)
	}
	if ConfigAdreses.AuditUrl != "" {
		obsURL := logger.NewObserverURL(ConfigAdreses.AuditUrl)
		PublisherAudit.RegisterObserver(obsURL)

	}
}
