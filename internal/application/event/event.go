package event

import (
	"context"
	"time"
)

// Event 代表領域事件的基礎介面
type Event interface {
	GetID() string
	GetType() string
	GetTimestamp() time.Time
	GetAggregateID() string
}

// BaseEvent 基礎事件結構
type BaseEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
	AggregateID string    `json:"aggregate_id"`
}

func (e BaseEvent) GetID() string           { return e.ID }
func (e BaseEvent) GetType() string         { return e.Type }
func (e BaseEvent) GetTimestamp() time.Time { return e.Timestamp }
func (e BaseEvent) GetAggregateID() string  { return e.AggregateID }

// EventHandler 事件處理器介面
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	CanHandle(eventType string) bool
}

// EventBus 事件總線介面
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
}

// EventStore 事件儲存介面
type EventStore interface {
	Save(ctx context.Context, event Event) error
	GetByAggregateID(ctx context.Context, aggregateID string) ([]Event, error)
	GetByType(ctx context.Context, eventType string) ([]Event, error)
}
