package routers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"labraboard/docs"
	"labraboard/internal/domains/iac/memory"
	"labraboard/internal/eventbus"
	api "labraboard/internal/routers/api"
	"labraboard/internal/services"
)

func InitRouter(publisher eventbus.EventPublisher, repository *memory.Repository) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.BestSpeed))

	iac, err := services.NewIacService(
		services.WithEventBus(publisher),
		services.WithRepository(repository))
	if err != nil {
		panic(err)
	}

	tfController, err := api.NewTerraformPlanController(iac)

	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", api.HelloWorld)
		}
		tf := v1.Group("/terraform/:projectId/plan/")
		{
			tf.GET("/:planId", tfController.GetTerraformPlan)
			tf.GET("/", tfController.FetchTerraformPlans)
			tf.POST("/", tfController.CreateTerraformPlan)
			tf.POST("/:planId/apply", tfController.ApplyTerraformPlan)
			tf.GET("/apply/:deploymentId", tfController.DeploymentTerraform)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
