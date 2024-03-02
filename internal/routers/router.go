package routers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"labraboard/docs"
	api "labraboard/internal/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.BestSpeed))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", api.HelloWorld)
		}
		tf := v1.Group("/terraform/:projectId/plan/")
		{
			tf.GET("/:planId", api.GetTerraformPlan)
			tf.GET("/", api.FetchTerraformPlans)
			tf.POST("/", api.CreateTerraformPlan)
			tf.POST("/:planId/apply", api.ApplyTerraformPlan)
			tf.GET("/apply/:deploymentId", api.DeploymentTerraform)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
