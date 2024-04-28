package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"labraboard"
	"labraboard/internal/eventbus/redisEventBus"
	"labraboard/internal/managers"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/postgres"
	"labraboard/internal/routers"
	"runtime"
)

func main() {
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
	ConfigRuntime()
	gin.SetMode(gin.ReleaseMode)
	db := postgres.NewDatabase(cfg.ConnectionString)
	defer func(db *postgres.Database) {
		err := db.Close()
		if err != nil {
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
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	eventBus, err := redisEventBus.NewRedisEventBus(context.Background(), redisEventBus.WithRedis(redisClient))
	if err != nil {
		panic(err)
	}
	//
	delayTaskManager, err := managers.NewDelayTaskManager(
		context.Background(),
		managers.WithRedis(redisClient),
		managers.WithEventPublisher(eventBus))

	if err != nil {
		panic(err)
	}

	//go ConfigureWorkers(eventBus, uow, delayTaskManager)
	routersInit := routers.InitRouter(eventBus, uow, delayTaskManager)
	err = routersInit.Run(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort))
	if err != nil {
		panic(err)
	}
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

//func ConfigureWorkers(subscriber eb.EventSubscriber, uow *repositories.UnitOfWork, mangerListener managers.DelayTaskMangerListener) {
//	go handlers.HandlePlan(subscriber, uow)
//}
