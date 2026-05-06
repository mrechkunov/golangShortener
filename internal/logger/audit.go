package logger

import (
	"fmt"
	"sync"

	"github.com/mrechkunov/golangShortener.git/internal/model"
)

// интерфейс публикатора
type Publisher interface {
	RegisterObserver(o Observer) // Регистрация наблюдателя
	RemoveObserver(o Observer)   // Удаление наблюдателя
	NotifyObservers()            // Уведомление всех наблюдателей
}

// интерфейс подписчиков (файл/url)
type Observer interface {
	Update(model.ObserverEvent) // Метод обновления
}

// реализвция publisher
type Audit struct {
	observers []Observer
	event     model.ObserverEvent
}

func (a *Audit) RegisterObserver(o Observer) {
	a.observers = append(a.observers, o)
}

func (a *Audit) RemoveObserver(o Observer) {
	for i, observer := range a.observers {
		if observer == o {
			a.observers = append(a.observers[:i], a.observers[i+1:]...)
			break
		}
	}
}

func (a *Audit) NotifyObservers() {
	for _, observer := range a.observers {
		observer.Update(a.event)
	}
}

// принимает новое событие и оповещает всех подписчиков
func (a *Audit) Event(newEvent model.ObserverEvent) {
	a.event = newEvent
	a.NotifyObservers()
}

// реализуем структуры и методы подписчиков
type ObserverFile struct {
	fileName string
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
			obsFile = &ObserverFile{fileName: auditFileName}
		})
	return obsFile
}

func (of *ObserverFile) Update(model.ObserverEvent) {
	// логика записи в файл
	fmt.Println("write to file")
}

type ObserverURL struct {
	URLName string
}

func NewObserverURL(auditURLName string) *ObserverURL {
	onceURL.Do( // функция ниже выполнится только один раз
		func() {
			// инициализируем объект
			obsURL = &ObserverURL{URLName: auditURLName}
		})
	return obsURL
}

func (su *ObserverURL) Update(model.ObserverEvent) {
	// логика записи события в url POST запрос
	fmt.Println("write to URL")
}
