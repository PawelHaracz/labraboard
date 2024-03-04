package eventbus

type EventPublisher interface {
	Publish(key Events, event interface{})
}

type EventSubscriber interface {
	Subscribe(key Events) chan interface{}
	Unsubscribe(key Events, ch chan interface{})
}

type EventBus struct {
	EventPublisher
	EventSubscriber
}
