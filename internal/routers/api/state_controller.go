package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/domains/iac/memory"
	"log"
	"net/http"
)

type payload struct {
	Data    []byte
	MD5     []byte
	Version int `json:"version"`
}

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
	//var p payload
	aggregate, err := c.Repository.GetState(uuid.MustParse(projectId))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return
	}
	state, err := aggregate.Deserialize()

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return
	}
	//p.Data = state
	//if raw := context.Request.Header.Get("Content-MD5"); raw != "" {
	//	md5, err := base64.StdEncoding.DecodeString(raw)
	//	if err != nil {
	//		context.Error(fmt.Errorf("Failed to decode Content-MD5 '%s': %s", raw, err))
	//	}
	//
	//	p.MD5 = md5
	//} else {
	//	// Generate the MD5
	//	hash := md5.Sum(p.Data)
	//	p.MD5 = hash[:]
	//}
	//p.Version = 4
	context.JSON(http.StatusOK, state)
	return
	//context.Writer.Write(stateBytes)
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
	ref := context.Query("ref")
	log.Default().Printf("ref: %s", ref)
	projectId := context.Param("projectId")
	aggregate, err := c.Repository.GetState(uuid.MustParse(projectId))
	if err != nil {
		aggregate, err = aggregates.NewTerraformState(uuid.MustParse(projectId), make([]byte, 0))
		if err != nil {
			context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
			return
		}
		err = c.Repository.AddState(aggregate)
		if err != nil {
			context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save to storage"})
			return
		}
	}
	var state map[string]interface{}
	if err := json.NewDecoder(context.Request.Body).Decode(&state); err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save to storage"})
		return
	}
	body, _ := json.Marshal(state)

	aggregate.SetState(&body)
	err = c.Repository.UpdateState(aggregate)
	if err != nil {
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
