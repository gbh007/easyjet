package eventbus

import (
	"log/slog"
	"sync"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type EventQueue struct {
	projectEventHandlers   map[string]chan entity.ProjectEvent
	projectEventHandlersMu *sync.RWMutex
	appEventHandlers       map[string]chan entity.AppEvent
	appEventHandlersMu     *sync.RWMutex
	logger                 *slog.Logger
}

func New(logger *slog.Logger) *EventQueue {
	return &EventQueue{
		projectEventHandlers:   make(map[string]chan entity.ProjectEvent),
		projectEventHandlersMu: &sync.RWMutex{},
		appEventHandlers:       make(map[string]chan entity.AppEvent),
		appEventHandlersMu:     &sync.RWMutex{},
		logger:                 logger,
	}
}

func (eq *EventQueue) PublishProjectEvent(event entity.ProjectEvent) {
	eq.projectEventHandlersMu.RLock()
	defer eq.projectEventHandlersMu.RUnlock()

	for name, ch := range eq.projectEventHandlers {
		select {
		case ch <- event:
		default:
			eq.logger.Warn("overflow project pubsub chan", "name", name)
		}
	}
}

func (eq *EventQueue) SubscribeProjectEvent(name string, c int) <-chan entity.ProjectEvent {
	eq.projectEventHandlersMu.Lock()
	defer eq.projectEventHandlersMu.Unlock()

	ch, ok := eq.projectEventHandlers[name]
	if !ok {
		ch = make(chan entity.ProjectEvent, c)
		eq.projectEventHandlers[name] = ch
	}

	return ch
}

func (eq *EventQueue) PublishAppEvent(event entity.AppEvent) {
	eq.appEventHandlersMu.RLock()
	defer eq.appEventHandlersMu.RUnlock()

	for name, ch := range eq.appEventHandlers {
		select {
		case ch <- event:
		default:
			eq.logger.Warn("overflow app pubsub chan", "name", name)
		}
	}
}

func (eq *EventQueue) SubscribeAppEvent(name string, c int) <-chan entity.AppEvent {
	eq.appEventHandlersMu.Lock()
	defer eq.appEventHandlersMu.Unlock()

	ch, ok := eq.appEventHandlers[name]
	if !ok {
		ch = make(chan entity.AppEvent, c)
		eq.appEventHandlers[name] = ch
	}

	return ch
}

func (eq *EventQueue) Close() {
	eq.projectEventHandlersMu.Lock()
	defer eq.projectEventHandlersMu.Unlock()

	for k, ch := range eq.projectEventHandlers {
		close(ch)
		delete(eq.projectEventHandlers, k)
	}

	eq.appEventHandlersMu.Lock()
	defer eq.appEventHandlersMu.Unlock()

	for k, ch := range eq.appEventHandlers {
		close(ch)
		delete(eq.appEventHandlers, k)
	}
}
