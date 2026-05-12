package repository

import (
	"errors"
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

// NewSafeMap return ptr to new map with mutex
func NewSafeMap() *SafeMap {
	return &SafeMap{
		m:       make(map[string]model.Event),
		counter: 1,
	}
}

// SetData insert to map new row with shortURL, originalURL, UID
func (s *SafeMap) SetData(shortURL string, originalURL string, cookie string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	uid, err := cryptoauth.GetIDFromCookie(cookie)
	if err != nil {
		logger.Log.Warnln("can not Get ID from cookie while setdata in storage")
	}
	newEvent := model.Event{
		ID:          s.counter,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		Cookie:      cookie,
		UID:         uid,
		IsDeleted:   false,
	}
	if _, ok := s.m[shortURL]; !ok {
		s.m[shortURL] = newEvent
		s.counter++
	} else {
		return errors.New("409 Conflict")
	}
	return nil
}

// Return originalURL from map by shortURL if row is not exist, return isFound = false
func (s *SafeMap) GetData(shortURL string) (originalURL string, isFound bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if val, ok := s.m[shortURL]; !ok {
		return "", false
	} else {
		originalURL = val.OriginalURL
		return originalURL, true
	}
}

// IsDeleted return true if shortURL is mark as deleted
func (s *SafeMap) IsDeleted(shortURL string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[shortURL].IsDeleted
}

// IsCookieExist return true if cookie is exist in map
func (s *SafeMap) IsCookieExist(cookie string) bool {
	// перебор всей мапы и сравнение поля cookie
	s.mu.RLock()
	defer s.mu.RUnlock()
	returnValue := false
	for _, value := range s.m {
		if value.Cookie == cookie {
			returnValue = true
		}
	}
	return returnValue
}

// GetDataByUID return batch of URLs whitch user set.
func (s *SafeMap) GetDataByUID(uid uint32) []model.ResponseDataBatchByCookie {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.ResponseDataBatchByCookie
	for _, value := range s.m {
		if value.UID == uid {
			var add = model.ResponseDataBatchByCookie{
				OriginalURL: value.OriginalURL,
				ShortURL:    value.ShortURL,
			}
			result = append(result, add)
		}
	}
	return result
}

// IsCreator return true if user is creator of shortURL else false
func (s *SafeMap) IsCreator(shortURL string, cookie string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.m[shortURL].Cookie == cookie {
		return true
	}
	return false
}

func (s *SafeMap) Close() error {
	return nil
}

// SetIsDeleted  mark all shortURLs from slice as deleted
func (s *SafeMap) SetIsDeleted(shortURLs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, val := range shortURLs {
		tmpm := s.m[val]
		tmpm.IsDeleted = true
		s.m[val] = tmpm
	}
}
