package event

import (
	"context"
	"golang_dev_docker/internal/application/event"
	"sync"
)

// InMemoryEventBus 記憶體事件總線實現
type InMemoryEventBus struct {
	handlers map[string][]event.EventHandler
	mu       sync.RWMutex
}

// NewInMemoryEventBus 創建記憶體事件總線
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]event.EventHandler),
	}
}

// Publish 發布事件
func (bus *InMemoryEventBus) Publish(ctx context.Context, evt event.Event) error {
	bus.mu.RLock()
	handlers, exists := bus.handlers[evt.GetType()]
	bus.mu.RUnlock()

	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if handler.CanHandle(evt.GetType()) {
			go func(h event.EventHandler) {
				h.Handle(ctx, evt)
			}(handler)
		}
	}

	return nil
}

// Subscribe 訂閱事件
func (bus *InMemoryEventBus) Subscribe(eventType string, handler event.EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if _, exists := bus.handlers[eventType]; !exists {
		bus.handlers[eventType] = make([]event.EventHandler, 0)
	}

	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消訂閱事件
func (bus *InMemoryEventBus) Unsubscribe(eventType string, handler event.EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	handlers, exists := bus.handlers[eventType]
	if !exists {
		return nil
	}

	for i, h := range handlers {
		if h == handler {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}
