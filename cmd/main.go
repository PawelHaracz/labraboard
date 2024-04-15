package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"labraboard"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/redis"
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
			panic("cannot read config file")
		}
	} else {
		err := cleanenv.ReadConfig(*configFile, &cfg)
		if err != nil {
			panic("cannot read config file")
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

	eventBus := redis.NewRedisEventBus(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB, context.Background())
	go ConfigureWorkers(eventBus, uow)
	routersInit := routers.InitRouter(eventBus, uow, db)
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

func ConfigureWorkers(subscriber eb.EventSubscriber, uow *repositories.UnitOfWork) {
	HandlePlan(subscriber, uow)
}
