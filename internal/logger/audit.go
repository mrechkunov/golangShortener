package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/mrechkunov/golangShortener.git/internal/model"
)

// Publisher interface in Audit method
type Publisher interface {
	RegisterObserver(o Observer) // Регистрация наблюдателя
	RemoveObserver(o Observer)   // Удаление наблюдателя
	NotifyObservers()            // Уведомление всех наблюдателей
}

// Observer interface in Audit method (file/url)
type Observer interface {
	Update(model.ObserverEvent) // Метод обновления
}

// Audit publisher
type Audit struct {
	observers []Observer
	event     model.ObserverEvent
	mu        sync.RWMutex
}

// Register new Observer
func (a *Audit) RegisterObserver(o Observer) {
	a.mu.Lock()
	a.observers = append(a.observers, o)
	a.mu.Unlock()
}

// Remove Observer
func (a *Audit) RemoveObserver(o Observer) {
	a.mu.Lock()
	for i, observer := range a.observers {
		if observer == o {
			a.observers = append(a.observers[:i], a.observers[i+1:]...)
			break
		}
	}
	a.mu.Unlock()
}

// Notify all Observers
func (a *Audit) NotifyObservers() {
	for _, observer := range a.observers {
		observer.Update(a.event)
	}
}

// Get new Event and Notify all Observers
func (a *Audit) Event(newEvent model.ObserverEvent) {
	a.mu.Lock()
	a.event = newEvent
	a.NotifyObservers()
	a.mu.Unlock()
}

// реализуем структуры и методы подписчиков
type ObserverFile struct {
	fileName string
	file     *os.File
	writer   *bufio.Writer
	mu       sync.RWMutex
}

var (
	onceFile, onceURL sync.Once
	obsFile           *ObserverFile
	obsURL            *ObserverURL
)

func NewObserverFile(auditFileName string) *ObserverFile {
	onceFile.Do( // функция ниже выполнится только один раз
		func() {
			// инициализируем объект
			file, err := os.OpenFile(auditFileName, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				Log.Warnln("can not open file to audit logging")
			}
			obsFile = &ObserverFile{
				fileName: auditFileName,
				file:     file,
				writer:   bufio.NewWriter(file)}

		})
	return obsFile
}

func (of *ObserverFile) Update(AuditData model.ObserverEvent) {
	data, err := json.Marshal(&AuditData)
	if err != nil {
		Log.Infoln("error while marshaling json", err)
	}
	of.mu.Lock()
	// записываем событие в буфер
	if _, err := of.writer.Write(data); err != nil {
		Log.Infoln("error while write to buffer", err)
		of.mu.Unlock()
	}

	// добавляем перенос строки
	if err := of.writer.WriteByte('\n'); err != nil {
		Log.Infoln("error while add new string", err)
		of.mu.Unlock()
	}
	// записываем буфер в файл
	of.writer.Flush()
	of.mu.Unlock()
}

type ObserverURL struct {
	URLName string
}

// Return ptr to new url observer
func NewObserverURL(auditURLName string) *ObserverURL {
	onceURL.Do( // функция ниже выполнится только один раз
		func() {
			// инициализируем объект
			obsURL = &ObserverURL{URLName: auditURLName}
		})
	return obsURL
}

// Update method for URL Observer
func (su *ObserverURL) Update(AuditData model.ObserverEvent) {
	// Маршалинг json
	jsonData, err := json.Marshal(AuditData)
	if err != nil {
		Log.Infoln("error while marshaling json", err)
	}
	// отправка POST запроса
	resp, err := retryablehttp.Post(su.URLName, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		Log.Infoln(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
}
