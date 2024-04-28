package redisEventBus

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"labraboard/internal/eventbus/events"
	"log"
)

type EventBusConfiguration func(os *EventBus) error

type EventBus struct {
	redisClient *redis.Client
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

	return eb, nil
}

func WithRedis(redisClient *redis.Client) EventBusConfiguration {
	return func(os *EventBus) error {
		os.redisClient = redisClient
		return nil
	}
}

func (r *EventBus) Subscribe(key events.EventName, ctx context.Context) chan []byte {
	subscriber := r.redisClient.Subscribe(ctx, string(key))

	item := make(chan []byte)
	go func() {
		defer close(item) //check it
		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				// handle error, for example log it and return
				log.Println(err)
				return
			}

			item <- []byte(msg.Payload)
			fmt.Println("Received message from " + msg.Channel + " channel.")
		}
	}()

	return item
}

func (r *EventBus) Unsubscribe(key events.EventName, ch chan []byte, ctx context.Context) {

}

func (r *EventBus) Publish(key events.EventName, event events.Event, ctx context.Context) {
	if err := r.redisClient.Publish(ctx, string(key), event).Err(); err != nil {
		panic(err)
	}
}
