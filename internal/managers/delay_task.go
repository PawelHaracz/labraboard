package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	_ "github.com/redis/go-redis/v9"
	"labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"time"
)

const delayedList = "delayed"

type DelayTaskManagerConfiguration func(os *delayTask) error

type DelayTaskManagerListener interface {
	Listen(ctx context.Context)
}

type DelayTaskManagerPublisher interface {
	Publish(EventName events.EventName, Content events.Event, WaitTime time.Duration, ctx context.Context)
}

type DelayTaskManager interface {
	DelayTaskManagerListener
	DelayTaskManagerPublisher
}

type delayTask struct {
	client    *redis.Client
	publisher eventbus.EventPublisher
}

type task struct {
	EventName events.EventName
	Content   events.Event
	WaitTime  time.Duration
}

func NewDelayTaskManager(ctx context.Context, configs ...DelayTaskManagerConfiguration) (DelayTaskManager, error) {
	dt := &delayTask{}
	for _, cfg := range configs {
		if err := cfg(dt); err != nil {
			return nil, err
		}
	}

	if dt.client == nil {
		return nil, errors.New("Redis client is not set ")
	}
	if dt.publisher == nil {
		return nil, errors.New("Event publisher is not set ")
	}

	cmd := dt.client.Ping(ctx)
	_, cmdErr := cmd.Result()
	if cmdErr != nil {
		return nil, errors.New("Cannot ping redisEventBus using client")
	}

	return dt, nil
}

func WithRedis(redisClient *redis.Client) DelayTaskManagerConfiguration {
	return func(os *delayTask) error {
		os.client = redisClient
		return nil
	}
}

func WithEventPublisher(publisher eventbus.EventPublisher) DelayTaskManagerConfiguration {
	return func(os *delayTask) error {
		os.publisher = publisher
		return nil
	}
}

func (dt *delayTask) Listen(ctx context.Context) {
	log := logger.GetWitContext(ctx)
	maxTime := time.Now().Unix()
	opt := &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", 0),
		Max: fmt.Sprintf("%d", maxTime),
	}
	cmd := dt.client.ZRevRangeByScore(ctx, delayedList, opt)
	resultSet, err := cmd.Result()
	if err != nil {
		log.Error().Err(err)
	}

	tasks := make([]task, len(resultSet))

	if len(tasks) == 0 {
		log.Trace().Msg("nothing to publish")
		return
	}
	///todo fix it: {"EventName":"lease_lock","Content":{"Id":"dfd5677a-7b62-4444-8f81-5dfad0ce1627","Type":"terraform","LeaseTime":"2024-06-05T19:59:38.160871Z"},"WaitTime":3600000000000}
	for i, t := range resultSet {
		err = json.Unmarshal([]byte(t), &tasks[i])
		if err != nil {
			log.Error().Err(err).Msg("JSON!!!")
			return
		}
		eventName := tasks[i].EventName
		dt.publisher.Publish(eventName, tasks[i].Content, ctx)

	}

	_, err = dt.client.ZRem(ctx, delayedList, resultSet).Result()
	if err != nil {
		log.Error().Err(err).Msg("redis_error")
		return
	}

}

func (dt *delayTask) Publish(EventName events.EventName, Content events.Event, WaitTime time.Duration, ctx context.Context) {
	t := &task{
		EventName,
		Content,
		WaitTime,
	}
	log := logger.GetWitContext(ctx)
	jsonValue, err := json.Marshal(t)
	if err != nil {
		log.Error().Err(err).Msg("JSON!!!")
		return
	}
	taskReadyInSeconds := time.Now().Add(WaitTime).Unix()
	member := redis.Z{
		Score:  float64(taskReadyInSeconds),
		Member: jsonValue,
	}
	_, err = dt.client.ZAdd(log.WithContext(ctx), delayedList, member).Result()
	if err != nil {
		log.Error().Err(err)
	}
}
