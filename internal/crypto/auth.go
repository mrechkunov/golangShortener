package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

var key = "secret key"

// generate new UID
func GenerateUID() ([]byte, error) {
	id := make([]byte, 4)
	_, err := rand.Read(id)
	if err != nil {
		fmt.Println("error while generate UID", err)
		return id, err
	}
	return id, nil
}

// return signed UID
func SignedUID(uid []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write(uid)
	if err != nil {
		fmt.Println("error while sign UID", err)
		return "", err
	}
	sign := h.Sum(nil)
	result := append(uid, sign...)
	return hex.EncodeToString(result), err
}

// validate incoming UID
func ValidateIncomingUID(cookie string) (bool, error) {
	h := hmac.New(sha256.New, []byte(key))
	data, err := hex.DecodeString(cookie)
	if err != nil {
		fmt.Println("error while decoding incoming UID", err)
		return false, err
	}
	h.Write([]byte(data[:4]))
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[4:]) {
		return true, nil
	} else {
		return false, nil
	}
}

// get ID from cookie
func GetIDFromCookie(cookie string) (uint32, error) {
	data, err := hex.DecodeString(cookie)
	if err != nil {
		fmt.Println("error while decoding incoming coockies", err)
		return 1, err
	}
	return binary.BigEndian.Uint32(data[:4]), nil
}
