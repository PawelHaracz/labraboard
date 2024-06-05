package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"labraboard"
	"labraboard/internal/eventbus/redisEventBus"
	"labraboard/internal/handlers"
	"labraboard/internal/logger"
	"labraboard/internal/managers"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/postgres"
	"labraboard/internal/routers"
	"os"
	"os/signal"
	"runtime"
	"time"
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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	ctx := context.Background()
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	log.Info().Int("GOMAXPROCS", nuCPU).Msgf("Running with %d CPUs\n", nuCPU)

	gin.SetMode(gin.ReleaseMode)
	db := postgres.NewDatabase(cfg.ConnectionString)
	defer func(db *postgres.Database) {
		err := db.Close()
		if err != nil {
			log.Panic().Err(err)
		}
	}(db)
	db.Migrate()
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIaCRepositoryDbRepository(db),
		repositories.WithTerraformStateDbRepository(db),
		repositories.WithIacPlanRepositoryDbRepository(db),
	)
	if err != nil {
		log.Panic().Err(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	eventBus, err := redisEventBus.NewRedisEventBus(context.Background(), redisEventBus.WithRedis(redisClient))
	if err != nil {
		log.Panic().Err(err)
	}
	//
	delayTaskManager, err := managers.NewDelayTaskManager(
		ctx,
		managers.WithRedis(redisClient),
		managers.WithEventPublisher(eventBus))

	if err != nil {
		log.Panic().Err(err)
	}

	//go ConfigureWorkers(eventBus, uow, delayTaskManager)
	routersInit := routers.InitRouter(eventBus, uow, delayTaskManager)
	if err != nil {
		log.Panic().Err(err)
	}
	handlerFactory := handlers.NewEventHandlerFactory(eventBus, eventBus, uow, fmt.Sprintf("%s:%d", cfg.ServiceDiscovery, cfg.HttpPort))

	allHandlers, err := handlerFactory.RegisterAllHandlers()
	if err != nil {
		log.Panic().Err(err)
	}
	for _, handler := range append(allHandlers) {
		go handler.Handle(ctx)
	}

	log.Info().Int("httpPort", cfg.HttpPort).Msgf("started server on 0.0.0.0:%d", cfg.HttpPort)

	go func(ctx context.Context) {
		for {
			delayTaskManager.Listen(ctx)
			time.Sleep(1 * time.Minute)
		}
	}(ctx)

	err = routersInit.Run(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort))
}
