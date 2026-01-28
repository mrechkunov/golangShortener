package config

import (
	"flag"
	"fmt"
)

type str struct {
	ServerBindAdress   *string
	ResultServerAdress *string
}

func Init() {
	var str str
	str.ServerBindAdress = flag.String("adress", "localhost:8080", "adress to server run")
	flag.Parse()
	fmt.Println(&str.ServerBindAdress)
}
