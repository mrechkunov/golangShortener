package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mrechkunov/golangShortener.git/internal/cryptoauth"
	"github.com/mrechkunov/golangShortener.git/internal/logger"
	"github.com/mrechkunov/golangShortener.git/internal/model"
)

type SafeMap struct {
	mu      sync.RWMutex
	m       map[string]model.Event
	counter int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m:       make(map[string]model.Event),
		counter: 1,
	}
}

func (s *SafeMap) SetData(shortURL string, originalURL string, cookie string) error {
	s.mu.Lock() // Блокировка на запись
	defer s.mu.Unlock()
	uid, err := cryptoauth.GetIDFromCookie(cookie)
	if err != nil {
		logger.Log.Errorln("can not Get ID from cookie while setdata in storage")
	}
	newEvent := model.Event{
		ID:          s.counter,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		Cookie:      cookie,
		UID:         uid,
	}
	if _, ok := s.m[shortURL]; !ok {
		s.m[shortURL] = newEvent
		s.counter++
	} else {
		return errors.New("409 Conflict")
	}
	return nil
}

func (s *SafeMap) GetData(shortURL string) (string, bool) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	if val, ok := s.m[shortURL]; !ok {
		return "", false
	} else {
		originalURL := val.OriginalURL
		return originalURL, true
	}
}
func (s *SafeMap) IsCookieExist(cookie string) bool {
	// перебор всей мапы и сравнение поля cookie
	fmt.Println("try to exist in storage", cookie)
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()

	returnValue := false
	fmt.Println("try to find cookie")
	if len(s.m) > 0 {
		for _, value := range s.m {
			if value.Cookie == cookie {
				returnValue = true
			}
		}
	}
	return returnValue
}

func (s *SafeMap) Close() error {
	return nil
}
