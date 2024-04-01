package eventbus

import "context"

type EventPublisher interface {
	Publish(key EventName, event interface{}, ctx context.Context)
}

type EventSubscriber interface {
	Subscribe(key EventName, ctx context.Context) chan []byte
	Unsubscribe(key EventName, ch chan []byte, ctx context.Context)
}

type Bus struct {
	EventPublisher
	EventSubscriber
}
