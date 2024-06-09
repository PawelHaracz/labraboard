package routers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"labraboard/docs"
	"labraboard/internal/eventbus"
	"labraboard/internal/managers"
	"labraboard/internal/repositories"
	api "labraboard/internal/routers/api"
	"labraboard/internal/services"
)

func InitRouter(publisher eventbus.EventPublisher, unitOfWork *repositories.UnitOfWork, delayTaskManagerPublisher managers.DelayTaskManagerPublisher, frontendPath string) *gin.Engine {
	r := gin.New()

	r.Use(requestid.New(), UseCorrelationId(), GinLogger(), gin.Recovery(), gzip.Gzip(gzip.BestSpeed))
	r.Use(UnitedSetup(unitOfWork))

	r.Use(static.Serve("/", static.LocalFile(frontendPath, true)))

	iac, err := services.NewIacService(
		services.WithEventBus(publisher),
		services.WithUnitOfWork(unitOfWork),
		services.WithDelayTaskManagerPublisher(delayTaskManagerPublisher),
	)

	if err != nil {
		panic(err)
	}

	tfController, err := api.NewTerraformPlanController(iac)
	stateController, err := api.NewStateController(delayTaskManagerPublisher)
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
			project.PUT("/:projectId/env", iacController.AddEnv)
			project.DELETE("/:projectId/env/:envName", iacController.RemoveEnv)
			project.PUT("/:projectId/variable", iacController.AddVariable)
			project.DELETE("/:projectId/variable/:variableName", iacController.RemoveVariable)
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
			tf.POST("/:planId/schedule", tfController.SchedulePlan)
			tf.GET("/", tfController.FetchTerraformPlans)
			tf.POST("/", tfController.CreateTerraformPlan)
			tf.POST("/:planId/apply", tfController.ApplyTerraformPlan)
			tf.GET("/apply/:deploymentId", tfController.DeploymentTerraform)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
