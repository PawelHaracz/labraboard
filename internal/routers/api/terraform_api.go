package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateTerraformPlan
// @BasePath /api/v1
// @Summary Method to run Terraform Plan for a given project and return the plan id
// @Schemes http
// @Description
// @Tags terraform
// @Accept json
// @Produce json
// @Success 200 {string} CreateTerraformPlan
// @Router /terraform/plan [POST]
func CreateTerraformPlan(g *gin.Context) {

}

// GetTerraformPlan
// @BasePath /api/v1
// @Summary Method returns the terraform plan output for a given plan id
// @Schemes http
// @Param projectId path string true "project id"
// @Param planId path string true "plan id"
// @Description
// @Tags terraform
// @Accept json
// @Produce json
// @Success 200 {string} GetTerraformPlan
// @Router /terraform/{projectId}/plan/{planId} [GET]
func GetTerraformPlan(g *gin.Context) {
	planId := g.Param("planId")
	projectId := g.Param("projectId")
	g.String(http.StatusOK, "hello world %s %s", planId, projectId)
}

// ApplyTerraformPlan
// @BasePath /api/v1
// @Summary Method Apply changes for a given plan id return deployment id
// @Schemes http
// @Param projectId path string true "project id"
// @Param planId path string true "plan id"
// @Description do ping
// @Tags terraform
// @Accept json
// @Produce json
// @Success 200 {string} ApplyTerraformPlan
// @Router /terraform/{projectId}/plan/{planId}/apply [POST]
func ApplyTerraformPlan(g *gin.Context) {

}

// DeploymentTerraform
// @BasePath /api/v1
// @Summary Method to fetch deployment status for a given deployment id
// @Schemes http
// @Param projectId path string true "project id"
// @Param planId path string true "plan id"
// @Param deploymentId path string true "deployment id"
// @Description do ping
// @Tags terraform
// @Accept json
// @Produce json
// @Success 200 {string} DeploymentTerraform
// @Router /terraform/{projectId}/plan/{planId}/apply/{deploymentId} [GET]
func DeploymentTerraform(g *gin.Context) {

}

// FetchTerraformPlans
// @BasePath /api/v1
// @Summary Fetch all the terraform plans for a given project
// @Schemes http
// @Param projectId path string true "project id"
// @Description do ping
// @Tags terraform
// @Accept json
// @Produce json
// @Success 200 {string} FetchTerraformPlans
// @Router /terraform/{projectId}/plan [GET]
func FetchTerraformPlans(g *gin.Context) {

}

//https://github.com/lovemapa/todo-Go
