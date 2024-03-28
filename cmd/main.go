package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"labraboard"
	"labraboard/internal/eventbus/redis"
	dbmemory "labraboard/internal/repositories/memory"
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
	repository, err := dbmemory.NewRepository()
	if err != nil {
		panic(err)
	}
	go ConfigureWorkers(repository)

	eventBus := redis.NewRedisEventBus(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB, context.Background())
	routersInit := routers.InitRouter(eventBus, repository, db)
	err = routersInit.Run(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort))
	if err != nil {
		panic(err)
	}
	//https://www.squash.io/optimizing-gin-in-golang-project-structuring-error-handling-and-testing/
	//https://github.com/swaggo/gin-swagger
	//https://github.com/eddycjy/go-gin-example
	//https://github.com/derekahn/ultimate-go/blob/master/language/interfaces/main.go
	//https://github.com/percybolmer/ddd-go
	//https://velocity.tech/blog/build-a-microservice-based-application-in-golang-with-gin-redis-and-mongodb-and-deploy-it-in-k8s
	//https://www.ompluscator.com/article/golang/practical-ddd-domain-repository/?source=post_page-----d308c9d79ba7--------------------------------
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func ConfigureWorkers(repository *dbmemory.Repository) {
	handlePlan(repository)
}
