package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go/log"
	"labraboard/internal/logger"
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
	l := logger.GetWitContext(context)
	projects, err := iac.iac.GetProjects(l.WithContext(context))
	if err != nil {
		l.Warn().Err(err)
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
	l := logger.GetWitContext(context).
		With().
		Str("projectId", projectId).
		Logger()
	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		l.Warn().Err(err).Msg("cannot parsed uuid")
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	project, err := iac.iac.GetProject(parsedProjectId, l.WithContext(context))

	if err != nil {
		l.Warn().Err(err)
		log.Error(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve project"})
		return
	}

	url, branch, path := project.GetRepo()

	base := &dtos.GetProjectBaseDto{
		IacType: int(project.IacType),
		Id:      project.GetID(),
	}
	dto := &dtos.GetProjectDto{
		GetProjectBaseDto: *base,
		RepositoryUrl:     url,
		RepositoryBranch:  branch,
		TerraformPath:     path,
		Envs:              project.GetEnvs(true),
		Variables:         project.GetVariableMap(),
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
	l := logger.GetWitContext(context)

	if err := context.BindJSON(&dto); err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
	}
	var repo = &vo.IaCRepo{
		Url:           dto.RepositoryUrl,
		DefaultBranch: dto.RepositoryBranch,
		Path:          dto.TerraformPath,
	}

	id, err := iac.iac.CreateProject(vo.IaCType(dto.IacType), repo, l.WithContext(context))

	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve project"})
		return
	}
	context.JSON(http.StatusOK, id)
}

// AddEnv Add an env to a project
// @Summary Add new env to a project
// @Schemes
// @Description Add new env to a project
// @Tags project
// @Param projectId path string true "project id"
// @Param project body dtos.AddEnvDto true "Env to add"
// @Accept json
// @Produce json
// @Success 200 {string} interface{}
// @Router /project/{projectId}/env [PUT]
func (iac *IacController) AddEnv(context *gin.Context) {
	projectId := context.Param("projectId")
	l := logger.GetWitContext(context).
		With().
		Str("projectId", projectId).
		Logger()
	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		l.Warn().Err(err).Msg("cannot parsed uuid")
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	var dto dtos.AddEnvDto
	if err := context.BindJSON(&dto); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}
	err = iac.iac.AddEnv(parsedProjectId, dto.Name, dto.Value, dto.IsSecret, l.WithContext(context))
	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": fmt.Sprintf("cannot add env %s", dto.Name)})
		return

	}
	context.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("env %s added", dto.Name)})
}

// RemoveEnv Remove an env to a project
// @Summary Remove new env to a project
// @Schemes
// @Description Remove new env to a project
// @Tags project
// @Param projectId path string true "project id"
// @Param envName path string true "env name"
// @Accept json
// @Produce json
// @Success 200 {string} interface{}
// @Router /project/{projectId}/env/{envName} [DELETE]
func (iac *IacController) RemoveEnv(context *gin.Context) {
	projectId := context.Param("projectId")
	envName := context.Param("envName")

	l := logger.GetWitContext(context).
		With().
		Str("projectId", projectId).
		Str("envName", envName).
		Logger()

	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		l.Warn().Err(err).Msg("cannot parsed uuid")
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	if err = iac.iac.RemoveEnv(parsedProjectId, envName, l.WithContext(context)); err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": fmt.Sprintf("cannot remove env %s", envName)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("env %s removed", envName)})
}

// AddVariable Add a variable to a project
// @Summary Add new variable to a project
// @Schemes
// @Description Add new variable to a project
// @Tags project
// @Param projectId path string true "project id"
// @Param project body dtos.AddVariableDto true "Env to add"
// @Accept json
// @Produce json
// @Success 200 {string} interface{}
// @Router /project/{projectId}/variable [PUT]
func (iac *IacController) AddVariable(context *gin.Context) {
	projectId := context.Param("projectId")

	l := logger.GetWitContext(context).
		With().
		Str("projectId", projectId).
		Logger()

	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		l.Warn().Err(err).Msg("cannot parsed uuid")
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	var dto dtos.AddVariableDto
	if err = context.BindJSON(&dto); err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}
	err = iac.iac.AddVariable(parsedProjectId, dto.Name, dto.Value, l.WithContext(context))
	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": fmt.Sprintf("cannot add variable %s", dto.Name)})
		return

	}
	context.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("variable %s added", dto.Name)})
}

// RemoveVariable Remove a variable to a project
// @Summary Remove new variable to a project
// @Schemes
// @Description Remove new variable to a project
// @Tags project
// @Param projectId path string true "project id"
// @Param variableName path string true "variable name"
// @Accept json
// @Produce json
// @Success 200 {string} interface{}
// @Router /project/{projectId}/variable/{variableName} [DELETE]
func (iac *IacController) RemoveVariable(context *gin.Context) {
	projectId := context.Param("projectId")
	variableName := context.Param("variableName")

	l := logger.GetWitContext(context).
		With().
		Str("projectId", projectId).
		Str("variableName", variableName).
		Logger()

	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		l.Warn().Err(err).Msg("cannot parsed uuid")
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	if err = iac.iac.RemoveVariable(parsedProjectId, variableName, l.WithContext(context)); err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": fmt.Sprintf("cannot remove variable %s", variableName)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("variable %s removed", variableName)})
}
