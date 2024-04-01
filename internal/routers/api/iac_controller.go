package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"labraboard/internal/services"
	vo "labraboard/internal/valueobjects"
	"net/http"
)

type ProjectDto struct {
	IacType int `json:"type"`
}

type IacController struct {
	iac *services.IacService
}

func NewIacController(service *services.IacService) (*IacController, error) {
	return &IacController{iac: service}, nil
}

// GetProjects fetch all projects
// @Summary Get all projects
// @Schemes
// @Description projects
// @Tags project
// @Accept json
// @Produce json
// @Success 200 {array} aggregates.Iac
// @Router /project [GET]
func (iac *IacController) GetProjects(context *gin.Context) {
	projects, err := iac.iac.GetProjects()

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve projects"})
		return
	}

	context.JSON(http.StatusOK, projects) //TODO map to a struct
}

// GetProject Fetch a project by id
// @Summary Fetch a project by id
// @Schemes
// @Description Fetch a project by id
// @Param projectId path string true "project id"
// @Tags project
// @Accept json
// @Produce json
// @Success 200 {object} aggregates.Iac
// @Router /project/{projectId} [GET]
func (iac *IacController) GetProject(context *gin.Context) {
	projectId := context.Param("projectId")
	project, err := iac.iac.GetProject(uuid.MustParse(projectId))

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve project"})
		return
	}
	context.JSON(http.StatusOK, project)
}

// CreateProject Create a new project
// @Summary Create a new project
// @Schemes
// @Description Create a new project
// @Tags project
// @Param project body ProjectDto true "Create project"
// @Accept json
// @Produce json
// @Success 200 {string} projectId
// @Router /project [POST]
func (iac *IacController) CreateProject(context *gin.Context) {

	var dto ProjectDto
	if err := context.BindJSON(&dto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
	}

	id, err := iac.iac.CreateProject(vo.IaCType(dto.IacType))

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve project"})
		return
	}
	context.JSON(http.StatusOK, id)
}

//TODO IMPLEMENT DTOS AND REPOSITORY CONNECTION
