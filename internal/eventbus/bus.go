package eventbus

type EventPublisher interface {
	Publish(key EventName, event interface{})
}

type EventSubscriber interface {
	Subscribe(key EventName) chan interface{}
	Unsubscribe(key EventName, ch chan interface{})
}

type Bus struct {
	EventPublisher
	EventSubscriber
}
