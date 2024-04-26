package main

import (
	"context"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"labraboard"
	"labraboard/internal/eventbus/redis"
	"labraboard/internal/handlers"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/postgres"
	"os"
	"os/signal"
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

	go handlers.HandlePlan(eventBus, uow)
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}
