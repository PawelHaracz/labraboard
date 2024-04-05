package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"labraboard/internal/routers/api/dtos"
	"labraboard/internal/services"
	"net/http"
)

type TerraformPlanController struct {
	services.IacService
}

func NewTerraformPlanController(iac *services.IacService) (*TerraformPlanController, error) {
	return &TerraformPlanController{
		IacService: *iac,
	}, nil
}

// CreateTerraformPlan
// @BasePath /api/v1
// @Summary Method to run Terraform Plan for a given project and return the plan id
// @Schemes http
// @Param projectId path string true "project id"
// @Description
// @Tags terraform
// @Accept json
// @Produce json
// @Success 200 {string} CreateTerraformPlan
// @Router /terraform/{projectId}/plan [POST]
func (c *TerraformPlanController) CreateTerraformPlan(g *gin.Context) {
	projectId := g.Param("projectId")
	planId, err := c.IacService.RunTerraformPlan(uuid.MustParse(projectId))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	g.String(http.StatusOK, planId.String())

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
// @Success 200 {object} dtos.PlanWithOutputDto
// @Router /terraform/{projectId}/plan/{planId} [GET]
func (c *TerraformPlanController) GetTerraformPlan(g *gin.Context) {
	planId := g.Param("planId")
	projectId := g.Param("projectId")
	plan, err := c.IacService.GetPlan(uuid.MustParse(projectId), uuid.MustParse(planId))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, plan)
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
func (c *TerraformPlanController) ApplyTerraformPlan(g *gin.Context) {

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
func (c *TerraformPlanController) DeploymentTerraform(g *gin.Context) {

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
// @Success 200 {array} dtos.PlanDto
// @Router /terraform/{projectId}/plan [GET]
func (c *TerraformPlanController) FetchTerraformPlans(g *gin.Context) {
	projectId := g.Param("projectId")
	plans := c.IacService.GetPlans(uuid.MustParse(projectId))
	dto := make([]*dtos.PlanDto, 0)

	for _, plan := range plans {
		dto = append(dto, &dtos.PlanDto{Id: plan.Id.String(),
			Status:    string(plan.Status),
			CreatedOn: plan.CreatedOn,
		})
	}
	g.JSON(http.StatusOK, dto)
}

//TODO IMPLEMENT HANDLE ERROR ON DEFFER FOLDER AFTER CLEARING THE FOLDER
