package redisEventBus

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"sync"
)

type EventBusConfiguration func(os *EventBus) error

type EventBus struct {
	redisClient *redis.Client
	subs        map[events.EventName][]chan []byte
	mu          sync.RWMutex
}

func NewRedisEventBus(ctx context.Context, configs ...EventBusConfiguration) (*EventBus, error) {
	eb := &EventBus{}
	for _, cfg := range configs {
		if err := cfg(eb); err != nil {
			return nil, err
		}
	}

	if eb.redisClient == nil {
		return nil, errors.New("redisEventBus client is not set ")
	}

	cmd := eb.redisClient.Ping(ctx)
	_, cmdErr := cmd.Result()
	if cmdErr != nil {
		return nil, errors.New("Cannot ping redisEventBus using client")
	}

	eb.subs = make(map[events.EventName][]chan []byte)
	return eb, nil
}

func WithRedis(redisClient *redis.Client) EventBusConfiguration {
	return func(os *EventBus) error {
		os.redisClient = redisClient
		return nil
	}
}

func (r *EventBus) Subscribe(key events.EventName, ctx context.Context) chan []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	log := logger.GetWitContext(ctx)

	subscriber := r.redisClient.Subscribe(log.WithContext(ctx), string(key))

	channel := subscriber.Channel()

	ch := make(chan []byte, 1)
	r.subs[key] = append(r.subs[key], ch)
	go func() {
		for {
			select {
			case msg := <-channel:
				ch <- []byte(msg.Payload)
			}
		}
	}()
	return ch

}

func (r *EventBus) Unsubscribe(key events.EventName, ch chan []byte, ctx context.Context) {

}

func (r *EventBus) Publish(key events.EventName, event events.Event, ctx context.Context) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	log := logger.GetWitContext(ctx)
	if err := r.redisClient.Publish(ctx, string(key), event).Err(); err != nil {
		log.Error().Err(err)
		return
	}
}
