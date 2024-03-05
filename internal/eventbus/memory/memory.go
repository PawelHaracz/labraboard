package memory

import (
	eb "labraboard/internal/eventbus"
	"sync"
)

type PubSub struct {
	mu   sync.RWMutex
	subs map[eb.EventName][]chan interface{}
}

func (ps *PubSub) Subscribe(key eb.EventName) chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 1)
	ps.subs[key] = append(ps.subs[key], ch)

	return ch
}

func (ps *PubSub) Publish(key eb.EventName, event interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subs[key] {
		ch <- event
	}
}

func (ps *PubSub) Unsubscribe(key eb.EventName, ch chan interface{}) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, subscriber := range ps.subs[key] {
		if subscriber == ch {
			ps.subs[key] = append(ps.subs[key][:i], ps.subs[key][i+1:]...)
			close(ch)
			break
		}
	}
}

func newMemoryPublisher() *PubSub {
	return &PubSub{
		subs: make(map[eb.EventName][]chan interface{}),
	}
}

func NewMemoryEventBus() *eb.Bus {
	ps := newMemoryPublisher()

	return &eb.Bus{
		EventPublisher:  ps,
		EventSubscriber: ps,
	}
}
