package eventbus

import (
	"context"
	"labraboard/internal/eventbus/events"
)

type EventPublisher interface {
	Publish(key events.EventName, event events.Event, ctx context.Context)
}

type EventSubscriber interface {
	Subscribe(key events.EventName, ctx context.Context) chan []byte
	Unsubscribe(key events.EventName, ch chan []byte, ctx context.Context)
}

type Bus struct {
	EventPublisher
	EventSubscriber
}
