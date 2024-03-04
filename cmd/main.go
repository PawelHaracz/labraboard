package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	dbmemory "labraboard/internal/domains/iac/memory"
	eb "labraboard/internal/eventbus"
	ebmemory "labraboard/internal/eventbus/memory"
	"labraboard/internal/routers"
	"runtime"
)

var (
	eventBus = ebmemory.NewMemoryEventBus()
)

func main() {
	ConfigRuntime()
	go ConfigureWorkers()
	gin.SetMode(gin.ReleaseMode)
	repository, err := dbmemory.NewRepository()
	if err != nil {
		panic(err)
	}
	routersInit := routers.InitRouter(eventBus.EventPublisher, repository)
	err = routersInit.Run("0.0.0.0:8080")
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
			fmt.Println("Received message:", msg)
		}
	}()
}
