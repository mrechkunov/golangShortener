package logger

import "fmt"

func Audit(action string, userID uint32, url string) {
	fmt.Println(action, userID, url)
}
