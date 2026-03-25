package eventbus

import (
	"log/slog"
	"sync"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type EventQueue struct {
	eventHandlers   map[string]chan entity.Event
	eventHandlersMu *sync.RWMutex
	logger          *slog.Logger
}

func New(logger *slog.Logger) *EventQueue {
	return &EventQueue{
		eventHandlers:   make(map[string]chan entity.Event),
		eventHandlersMu: &sync.RWMutex{},
		logger:          logger,
	}
}

func (eq *EventQueue) PublishEvent(event entity.Event) {
	eq.eventHandlersMu.RLock()
	defer eq.eventHandlersMu.RUnlock()

	for name, ch := range eq.eventHandlers {
		select {
		case ch <- event:
		default:
			eq.logger.Warn("overflow project pubsub chan", "name", name)
		}
	}
}

func (eq *EventQueue) SubscribeEvent(name string, c int) <-chan entity.Event {
	eq.eventHandlersMu.Lock()
	defer eq.eventHandlersMu.Unlock()

	ch, ok := eq.eventHandlers[name]
	if !ok {
		ch = make(chan entity.Event, c)
		eq.eventHandlers[name] = ch
	}

	return ch
}

func (eq *EventQueue) Close() {
	eq.eventHandlersMu.Lock()
	defer eq.eventHandlersMu.Unlock()

	for k, ch := range eq.eventHandlers {
		close(ch)
		delete(eq.eventHandlers, k)
	}
}
