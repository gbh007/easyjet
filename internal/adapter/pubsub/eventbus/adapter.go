package eventbus

import (
	"log/slog"
	"sync"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type EventQueue struct {
	projectEventHandlers map[string]chan entity.ProjectEvent
	mu                   *sync.RWMutex
	logger               *slog.Logger
}

func New(logger *slog.Logger) *EventQueue {
	return &EventQueue{
		projectEventHandlers: make(map[string]chan entity.ProjectEvent),
		mu:                   &sync.RWMutex{},
		logger:               logger,
	}
}

func (eq *EventQueue) PublishProjectEvent(event entity.ProjectEvent) {
	eq.mu.RLock()
	defer eq.mu.RUnlock()

	for name, ch := range eq.projectEventHandlers {
		select {
		case ch <- event:
		default:
			eq.logger.Warn("overflow project pubsub chan", "name", name)
		}
	}
}

func (eq *EventQueue) SubscribeProjectEvent(name string, c int) <-chan entity.ProjectEvent {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	ch, ok := eq.projectEventHandlers[name]
	if !ok {
		ch = make(chan entity.ProjectEvent, c)
		eq.projectEventHandlers[name] = ch
	}

	return ch
}

func (eq *EventQueue) Close() {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	for k, ch := range eq.projectEventHandlers {
		close(ch)
		delete(eq.projectEventHandlers, k)
	}
}
