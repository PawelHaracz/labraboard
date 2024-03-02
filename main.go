package main

import (
	"github.com/gin-gonic/gin"
	"labraboard/internal/routers"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	routersInit := routers.InitRouter()
	err := routersInit.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	//https://www.squash.io/optimizing-gin-in-golang-project-structuring-error-handling-and-testing/
	//https://github.com/swaggo/gin-swagger
	//https://github.com/eddycjy/go-gin-example
}

//func
