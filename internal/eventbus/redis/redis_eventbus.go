package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"labraboard/internal/eventbus"
	"log"
)

type RedisEventBus struct {
	redisClient *redis.Client
}

func NewRedisEventBus(host string, port int, password string, db int, ctx context.Context) *RedisEventBus {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})
	redisClient.Ping(ctx)
	return &RedisEventBus{
		redisClient,
	}
}

func (r *RedisEventBus) Subscribe(key eventbus.EventName, ctx context.Context) chan interface{} {
	subscriber := r.redisClient.Subscribe(ctx, string(key))

	item := make(chan interface{})
	go func() {
		defer close(item) //checkit
		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				// handle error, for example log it and return
				log.Println(err)
				return
			}

			var payload interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				// handle error, for example log it and return
				log.Println(err)
				return
			}

			item <- payload
			fmt.Println("Received message from " + msg.Channel + " channel.")
		}
	}()

	return item
}

func (r *RedisEventBus) Unsubscribe(key eventbus.EventName, ch chan interface{}, ctx context.Context) {

}

func (r *RedisEventBus) Publish(key eventbus.EventName, event interface{}, ctx context.Context) {
	if err := r.redisClient.Publish(ctx, "send-user-data", event).Err(); err != nil {
		panic(err)
	}
}
