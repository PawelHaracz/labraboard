package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"labraboard/internal/domains/iac/memory"
	"net/http"
)

// https://expeditor.chef.io/docs/getting-started/terraform/
// https://github.com/platformod/united/blob/main/handlers.go
type StateController struct {
	*memory.Repository
}

func NewStateController(repository *memory.Repository) (*StateController, error) {
	return &StateController{
		Repository: repository}, nil
}

// GetState
// @BasePath /api/v1
// @Summary Get terraform state for a given project
// @Schemes http
// @Param projectId path string true "project id"
// @Description
// @Tags terraform
// @Tags state
// @Accept json
// @Produce json
// @Router /state/terraform/{projectId} [GET]
func (c *StateController) GetState(context *gin.Context) {
	projectId := context.Param("projectId")
	aggregate, err := c.Repository.GetState(uuid.MustParse(projectId))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}
	state, err := aggregate.GetState()
	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return
	}
	context.JSON(http.StatusOK, state)
}

// UpdateState
// @BasePath /api/v1
// @Summary Update terraform state for a given project
// @Schemes http
// @Param projectId path string true "project id"
// @Description
// @Tags terraform
// @Tags state
// @Accept json
// @Produce json
// @Router /state/terraform/{projectId} [POST]
func (c *StateController) UpdateState(context *gin.Context) {
	projectId := context.Param("projectId")
	aggregate, err := c.Repository.GetState(uuid.MustParse(projectId))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}
	body, _ := io.ReadAll(context.Request.Body)
	aggregate.SetState(body)
	if err := c.Repository.AddState(aggregate); err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save to storage"})
		return
	}
}

// Lock
// @BasePath /api/v1
// @Summary Lock terraform state for a given project
// @Schemes http
// @Param projectId path string true "project id"
// @Description
// @Tags terraform
// @Tags state
// @Accept json
// @Produce json
// @Success 200
func (c *StateController) Lock(context *gin.Context) {

}

// Unlock
// @BasePath /api/v1
// @Summary Unlock terraform state for a given project
// @Schemes http
// @Param projectId path string true "project id"
// @Description
// @Tags terraform
// @Tags state
// @Accept json
// @Produce json
// @Success 200
func (c *StateController) Unlock(context *gin.Context) {

}
