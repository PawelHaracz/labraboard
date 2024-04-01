package memory

import (
	"context"
	"encoding/json"
	eb "labraboard/internal/eventbus"
	"sync"
)

type PubSub struct {
	mu   sync.RWMutex
	subs map[eb.EventName][]chan []byte
}

func (ps *PubSub) Subscribe(key eb.EventName, ctx context.Context) chan []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan []byte, 1)
	ps.subs[key] = append(ps.subs[key], ch)

	return ch
}

func (ps *PubSub) Publish(key eb.EventName, event interface{}, ctx context.Context) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subs[key] {
		b, _ := json.Marshal(event)
		ch <- b
	}
}

func (ps *PubSub) Unsubscribe(key eb.EventName, ch chan []byte, ctx context.Context) {
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
		subs: make(map[eb.EventName][]chan []byte),
	}
}

func NewMemoryEventBus() *eb.Bus {
	ps := newMemoryPublisher()

	return &eb.Bus{
		EventPublisher:  ps,
		EventSubscriber: ps,
	}
}
