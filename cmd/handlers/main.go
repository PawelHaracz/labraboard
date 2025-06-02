package main

import (
	"context"
	"flag"
	"fmt"
	"labraboard"
	"labraboard/internal/eventbus/redisEventBus"
	"labraboard/internal/handlers"
	"labraboard/internal/logger"
	"labraboard/internal/managers"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/postgres"
	"os"
	"os/signal"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var cfg labraboard.Config
var log zerolog.Logger

func init() {
	configFile := flag.String("config", "", "config file, if empty then use env variables")
	flag.Parse()
	if *configFile == "" {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			panic(errors.Wrap(err, "Missing env variables"))
		}
	} else {
		err := cleanenv.ReadConfig(*configFile, &cfg)
		if err != nil {
			panic(errors.Wrap(err, "cannot read config file"))
		}
	}
	logger.Init(cfg.LogLevel, cfg.UsePrettyLogs)
	log = logger.Get()
}

func main() {
	log.Info().Msg("Starting handlers")
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())

	db := postgres.NewDatabase(cfg.ConnectionString)
	defer func(db *postgres.Database) {
		err := db.Close()
		if err != nil {
			cancel()
			log.Panic().Err(err)
		}
	}(db)
	db.Migrate()
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIaCRepositoryDbRepository(db),
		repositories.WithTerraformStateDbRepository(db),
		repositories.WithIacPlanRepositoryDbRepository(db),
		repositories.WithIacDeploymentRepositoryDbRepository(db),
	)
	if err != nil {
		cancel()
		log.Panic().Err(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	eventBus, err := redisEventBus.NewRedisEventBus(ctx, redisEventBus.WithRedis(redisClient))
	if err != nil {
		cancel()
		log.Panic().Err(err)
	}

	handlerFactory := handlers.NewEventHandlerFactory(eventBus, eventBus, uow, fmt.Sprintf("%s:%d", cfg.ServiceDiscovery, cfg.HttpPort))

	allHandlers, err := handlerFactory.RegisterAllHandlers()
	if err != nil {
		cancel()
		log.Panic().Err(err)
	}
	for _, handler := range allHandlers {
		go handler.Handle(ctx)
	}

	delayTaskManager, err := managers.NewDelayTaskManager(ctx,
		managers.WithRedis(redisClient),
		managers.WithEventPublisher(eventBus))

	if err != nil {
		cancel()
		log.Panic().Err(err)
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
