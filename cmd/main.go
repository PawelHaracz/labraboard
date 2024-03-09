package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"labraboard"
	dbmemory "labraboard/internal/domains/iac/memory"
	"labraboard/internal/domains/iac/postgres"
	eb "labraboard/internal/eventbus"
	ebmemory "labraboard/internal/eventbus/memory"
	"labraboard/internal/routers"
	iacSvc "labraboard/internal/services/iac"
	"runtime"
)

var (
	eventBus = ebmemory.NewMemoryEventBus()
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
	go ConfigureWorkers()
	gin.SetMode(gin.ReleaseMode)
	db := postgres.NewDatabase(cfg.ConnectionString)
	defer func(db *postgres.Database) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	//todo add repository for postgresql
	repository, err := dbmemory.NewRepository()
	if err != nil {
		panic(err)
	}
	routersInit := routers.InitRouter(eventBus.EventPublisher, repository)
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

func ConfigureWorkers() {
	pl := eventBus.Subscribe(eb.TRIGGERED_PLAN)
	//defer eventBus.Unsubscribe(eb.TRIGGERED_PLAN, pl)

	go func() {
		for msg := range pl {
			switch obj := msg.(type) {
			case uuid.UUID:
				fmt.Println("Received message:", msg)
				tofu, err := iacSvc.NewTofuIacService("")
				if err != nil {
					fmt.Println("error:", err)
				}
				_, err = tofu.Plan(obj)
				if err != nil {
					panic(err)
				}
			default:
				fmt.Errorf("cannot handle message type %T", obj)
			}

		}
	}()
}
