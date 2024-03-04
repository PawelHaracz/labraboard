package memory

import (
	eb "labraboard/internal/eventbus"
	"sync"
)

type PubSub struct {
	mu   sync.RWMutex
	subs map[eb.Events][]chan interface{}
}

func (ps *PubSub) Subscribe(key eb.Events) chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 1)
	ps.subs[key] = append(ps.subs[key], ch)

	return ch
}

func (ps *PubSub) Publish(key eb.Events, event interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subs[key] {
		ch <- event
	}
}

func (ps *PubSub) Unsubscribe(key eb.Events, ch chan interface{}) {
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
		subs: make(map[eb.Events][]chan interface{}),
	}
}

func NewMemoryEventBus() *eb.EventBus {
	ps := newMemoryPublisher()

	return &eb.EventBus{
		EventPublisher:  ps,
		EventSubscriber: ps,
	}
}
