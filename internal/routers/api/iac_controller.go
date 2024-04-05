package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"labraboard/internal/routers/api/dtos"
	"labraboard/internal/services"
	vo "labraboard/internal/valueobjects"
	"net/http"
)

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
// @Success 200 {array} dtos.GetProjectBaseDto
// @Router /project [GET]
func (iac *IacController) GetProjects(context *gin.Context) {
	projects, err := iac.iac.GetProjects()

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve projects"})
		return
	}
	var projectsDto = make([]*dtos.GetProjectBaseDto, 0)
	for _, project := range projects {
		if project == nil {
			continue
		}
		projectsDto = append(projectsDto, &dtos.GetProjectBaseDto{
			IacType: int(project.IacType),
			Id:      project.GetID(),
		})
	}
	context.JSON(http.StatusOK, projectsDto)
}

// GetProject Fetch a project by id
// @Summary Fetch a project by id
// @Schemes
// @Description Fetch a project by id
// @Param projectId path string true "project id"
// @Tags project
// @Accept json
// @Produce json
// @Success 200 {object} dtos.GetProjectDto
// @Router /project/{projectId} [GET]
func (iac *IacController) GetProject(context *gin.Context) {
	projectId := context.Param("projectId")
	project, err := iac.iac.GetProject(uuid.MustParse(projectId))

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve project"})
		return
	}
	base := &dtos.GetProjectBaseDto{
		IacType: int(project.IacType),
		Id:      project.GetID(),
	}
	dto := &dtos.GetProjectDto{
		GetProjectBaseDto: *base,
		RepositoryUrl:     project.Repo.Url,
		RepositoryBranch:  project.Repo.DefaultBranch,
		TerraformPath:     project.Repo.Path,
	}
	context.JSON(http.StatusOK, dto)
}

// CreateProject Create a new project
// @Summary Create a new project
// @Schemes
// @Description Create a new project
// @Tags project
// @Param project body dtos.CreateProjectDto true "Create project"
// @Accept json
// @Produce json
// @Success 200 {string} projectId
// @Router /project [POST]
func (iac *IacController) CreateProject(context *gin.Context) {

	var dto dtos.CreateProjectDto
	if err := context.BindJSON(&dto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
	}
	var repo = &vo.IaCRepo{
		Url:           dto.RepositoryUrl,
		DefaultBranch: dto.RepositoryBranch,
		Path:          dto.TerraformPath,
	}

	id, err := iac.iac.CreateProject(vo.IaCType(dto.IacType), repo)

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve project"})
		return
	}
	context.JSON(http.StatusOK, id)
}
