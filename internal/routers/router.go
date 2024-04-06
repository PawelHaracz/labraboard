package routers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"labraboard/docs"
	"labraboard/internal/eventbus"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/postgres"
	api "labraboard/internal/routers/api"
	"labraboard/internal/services"
)

func InitRouter(publisher eventbus.EventPublisher, unitOfWork *repositories.UnitOfWork, database *postgres.Database) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.BestSpeed))
	r.Use(UnitedSetup(unitOfWork))

	iac, err := services.NewIacService(
		services.WithEventBus(publisher),
		services.WithUnitOfWork(unitOfWork))
	if err != nil {
		panic(err)
	}

	tfController, err := api.NewTerraformPlanController(iac)
	stateController, err := api.NewStateController()
	iacController, err := api.NewIacController(iac)

	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group(docs.SwaggerInfo.BasePath)
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", api.HelloWorld)
		}
		project := v1.Group("/project")
		{
			project.GET("/", iacController.GetProjects)
			project.POST("/", iacController.CreateProject)
			project.GET("/:projectId", iacController.GetProject)
		}
		state := v1.Group("/state/terraform")
		{
			state.GET(":projectId", stateController.GetState)
			state.POST(":projectId", stateController.UpdateState)
			state.Handle("LOCK", ":projectId/lock", stateController.Lock)
			state.Handle("UNLOCK", ":projectId/lock", stateController.Unlock)
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
