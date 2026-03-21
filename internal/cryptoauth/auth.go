package cryptoauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mrechkunov/golangShortener.git/internal/logger"
)

var secretKey = "secret key"

// generate NEW UID, sign it and return cookieString
func GenerateNewCookie() string {
	id := make([]byte, 4)
	_, err := rand.Read(id)
	if err != nil {
		fmt.Println("error while generate UID", err)
	}
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(id)
	if err != nil {
		fmt.Println("error while sign UID", err)
	}
	sign := h.Sum(nil)
	result := append(id, sign...)
	return hex.EncodeToString(result)
}

// return signed UID
// func SignUID(uid []byte) (string, error) {
// 	h := hmac.New(sha256.New, []byte(secretKey))
// 	_, err := h.Write(uid)
// 	if err != nil {
// 		fmt.Println("error while sign UID", err)
// 		return "", err
// 	}
// 	sign := h.Sum(nil)
// 	result := append(uid, sign...)
// 	return hex.EncodeToString(result), err
// }

// validate cookie signature
func ValidateCookieSign(cookie string) (bool, error) {
	h := hmac.New(sha256.New, []byte(secretKey))
	data, err := hex.DecodeString(cookie)
	if err != nil {
		logger.Log.Warnln("error while decoding incoming cookie", err)
		return false, err
	}
	h.Write([]byte(data[:4]))
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[4:]) {
		return true, nil
	} else {
		return false, errors.New("not valid cookie signature")
	}
}

// get ID from cookie
func GetIDFromCookie(cookie string) (uint32, error) {
	data, err := hex.DecodeString(cookie)
	if err != nil {
		fmt.Println("error while decoding incoming cookie", err)
		return 0, err
	}
	return binary.BigEndian.Uint32(data[:4]), nil
}
