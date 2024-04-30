package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"labraboard"
	"labraboard/internal/eventbus/redisEventBus"
	"labraboard/internal/handlers"
	"labraboard/internal/managers"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/postgres"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	var cfg labraboard.Config
	configFile := flag.String("config", "", "config file, if empty then use env variables")
	flag.Parse()
	if *configFile == "" {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			panic(errors.Wrap(err, "cannot read config file"))
		}
	} else {
		err := cleanenv.ReadConfig(*configFile, &cfg)
		if err != nil {
			panic(errors.Wrap(err, "cannot read config file"))
		}
	}

	db := postgres.NewDatabase(cfg.ConnectionString)
	defer func(db *postgres.Database) {
		err := db.Close()
		if err != nil {
			cancel()
			panic(err)
		}
	}(db)
	db.Migrate()
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIaCRepositoryDbRepository(db),
		repositories.WithTerraformStateDbRepository(db),
		repositories.WithIacPlanRepositoryDbRepository(db),
	)
	if err != nil {
		cancel()
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	eventBus, err := redisEventBus.NewRedisEventBus(ctx, redisEventBus.WithRedis(redisClient))
	if err != nil {
		cancel()
		panic(err)
	}

	handlerFactory := handlers.NewEventHandlerFactory(eventBus, uow)

	allHandlers, err := handlerFactory.RegisterAllHandlers()
	if err != nil {
		cancel()
		panic(err)
	}
	for _, handler := range append(allHandlers) {
		go handler.Handle(ctx)
	}

	delayTaskManager, err := managers.NewDelayTaskManager(ctx,
		managers.WithRedis(redisClient),
		managers.WithEventPublisher(eventBus))

	if err != nil {
		cancel()
		panic(err)
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				return
			default:
				delayTaskManager.Listen(ctx)
			}
			time.Sleep(1 * time.Minute)
		}
	}(ctx)

	<-signalChan
	cancel()
}
