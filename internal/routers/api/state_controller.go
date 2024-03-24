package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/domains/iac/memory"
	"log"
	"net/http"
	"time"
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
	if state == nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	context.JSON(http.StatusOK, state)
	return
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
	id := context.Query("ID")
	projectId := context.Param("projectId")
	aggregate, err := c.Repository.GetState(uuid.MustParse(projectId))
	if err != nil {
		utc := time.Now().UTC()
		aggregate, err = aggregates.NewTerraformState(uuid.MustParse(projectId), make([]byte, 0), utc, utc, make([]byte, 0))
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
	storedLock, err := aggregate.GetLockInfo()
	if id != "" && storedLock.ID != id {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Locked by different ID", "ID": storedLock.ID})
		return
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
	projectId := context.Param("projectId")

	var reqLock aggregates.LockInfo
	_ = context.BindJSON(&reqLock)
	aggregate, err := c.Repository.GetState(uuid.MustParse(projectId))
	if err != nil {
		utc := time.Now().UTC()
		aggregate, err = aggregates.NewTerraformState(uuid.MustParse(projectId), make([]byte, 0), utc, utc, make([]byte, 0))
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

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return

	}

	storedLock, err := aggregate.GetLockInfo()
	if err != nil {
		context.JSON(http.StatusLocked, gin.H{"message": "Already Locked"})
		return
	}
	if storedLock != nil {
		context.JSON(http.StatusLocked, gin.H{"message": "Already Locked", "ID": storedLock.ID})
		return
	}
	//todo handle lock time
	if err := aggregate.SetLockInfo(&reqLock); err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save to storage"})
		return
	}
	context.JSON(http.StatusOK, reqLock.ID)
	return

	//
	//rc := c.MustGet("rc").(*redis.Client)
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
	projectId := context.Param("projectId")
	var reqLock aggregates.LockInfo
	_ = context.BindJSON(&reqLock)
	aggregate, err := c.Repository.GetState(uuid.MustParse(projectId))
	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return
	}

	err = aggregate.LeaseLock(&reqLock)

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}
