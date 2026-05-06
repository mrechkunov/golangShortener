package logger

import (
	"fmt"
	"time"

	"github.com/mrechkunov/golangShortener.git/internal/model"
)

type Publisher interface {
	register(Observer, string)
	deregister(Observer)
	notify(model.ObserverEvent)
}

// подписчики (файл/url)
type Observer interface {
	update(event model.ObserverEvent)
}

// реализвция publisher
type Event struct {
	Observers   []Observer
	Description string
}

func (e *Event) register(o Observer, descript string) {
	e.Observers = append(e.Observers, o)
	e.Description = descript
}
func (e *Event) deregister(o Observer) {
	for i, observer := range e.Observers {
		if observer == o {
			e.Observers = append(e.Observers[:i], e.Observers[i+1:]...)
			break
		}
	}
}
func (e *Event) notify(newEvent model.ObserverEvent) {
	for _, observer := range e.Observers {
		observer.update(newEvent)
	}
}

func Audit(action string, userID uint32, url string) {

	event := model.ObserverEvent{
		Ts:          time.Now(),
		Action:      action,
		UserId:      userID,
		OriginalURL: url,
	}
	fmt.Println(event)
}
