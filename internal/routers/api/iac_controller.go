package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"labraboard/internal/services"
	"net/http"
)

type IacController struct {
	iac *services.IacService
}

func NewIacController(service *services.IacService) (*IacController, error) {
	return &IacController{iac: service}, nil
}

func (iac *IacController) GetProjects(context *gin.Context) {
	projects, err := iac.iac.GetProjects()

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve projects"})
		return
	}

	context.JSON(http.StatusOK, projects) //TODO map to a struct
}

func (iac *IacController) GetProject(context *gin.Context) {
	projectId := context.Param("projectId")
	project, err := iac.iac.GetProject(uuid.MustParse(projectId))

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "cannot retrieve project"})
		return
	}
	context.JSON(http.StatusOK, project)
}

//func (iac *IacController) CreateProject(context *gin.Context) {
//	project, err := iac.iac.CreateProject()
//}
