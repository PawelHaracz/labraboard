package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/helpers"
	"labraboard/internal/logger"
	"labraboard/internal/managers"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
	"net/http"
	"time"
)

type StateController struct {
	delayTaskManger managers.DelayTaskManagerPublisher
}

func NewStateController(delayTaskManger managers.DelayTaskManagerPublisher) (*StateController, error) {
	return &StateController{delayTaskManger}, nil
}

// GetState
// @BasePath /api/v1
// @Summary Get terraform state for a given project
// @Schemes http
// @Param projectId path string true "project id"
// @Description
// @Tags state
// @Accept json
// @Produce json
// @Router /state/terraform/{projectId} [GET]
func (c *StateController) GetState(context *gin.Context) {
	projectId := context.Param("projectId")
	l := logger.GetWitContext(context).
		With().
		Str("projectId", projectId).
		Logger()

	repo := context.MustGet(string(helpers.UnitOfWorkSetup)).(*repositories.UnitOfWork).TerraformStateDbRepository
	parsedGuid, err := uuid.Parse(projectId)
	if err != nil {
		l.Warn().Err(err).Msg("cannot parsed uuid")
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	aggregate, err := repo.Get(parsedGuid)
	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	state, err := aggregate.Deserialize()
	if err != nil {
		l.Warn().Err(err).Msg("Cannot deserialize state")
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return
	}
	if state == nil {
		l.Warn().Msg("Not found state")
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
// @Tags state
// @Accept json
// @Produce json
// @Router /state/terraform/{projectId} [POST]
func (c *StateController) UpdateState(context *gin.Context) {
	ref := context.Query("ref")
	id := context.Query("ID")
	projectId := context.Param("projectId")
	l := logger.GetWitContext(context).With().Str("ref", ref).Str("storedLockId", id).Str("projectId", projectId).Logger()
	repo := context.MustGet(string(helpers.UnitOfWorkSetup)).(*repositories.UnitOfWork).TerraformStateDbRepository

	parsedGuid, err := uuid.Parse(projectId)
	if err != nil {
		l.Warn().Err(err).Msg("cannot parsed uuid")
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}
	aggregate, err := repo.Get(parsedGuid)
	if err != nil {
		l.Info().Err(err).Msg("terraform state doesn't exist, creating")
		utc := time.Now().UTC()
		aggregate, err = aggregates.NewTerraformState(uuid.MustParse(projectId), make([]byte, 0), utc, utc, make([]byte, 0))
		if err != nil {
			l.Warn().Err(err).Msg("cannot create state aggregate")
			context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
			return
		}
		err = repo.Add(aggregate)
		if err != nil {
			l.Warn().Err(err).Msg("Could not save to storage")
			context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save to storage"})
			return
		}
	}
	storedLock, err := aggregate.GetLockInfo()
	if id != "" && storedLock.ID != id {
		l.Warn().Msgf("Locked by different ID: %s", storedLock.ID)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Locked by different ID", "ID": storedLock.ID})
		return
	}
	var state map[string]interface{}
	if err = json.NewDecoder(context.Request.Body).Decode(&state); err != nil {
		l.Warn().Err(err).Msg("Could not deserialize state")
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not deserialize state"})
		return
	}
	body, _ := json.Marshal(state)

	aggregate.SetState(&body)
	err = repo.Update(aggregate)
	if err != nil {
		l.Warn().Err(err).Msg("Could not update to storage")
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
// @Tags state
// @Accept json
// @Produce json
// @Success 200
func (c *StateController) Lock(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	l := logger.GetWitContext(ctx).
		With().
		Str("projectId", projectId).
		Logger()

	var reqLock aggregates.LockInfo
	err := ctx.BindJSON(&reqLock)
	if err != nil {
		l.Warn().Err(err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Cannot bind object"})
		return
	}
	repo := ctx.MustGet(string(helpers.UnitOfWorkSetup)).(*repositories.UnitOfWork).TerraformStateDbRepository
	aggregate, err := repo.Get(uuid.MustParse(projectId))
	if err != nil {
		l.Info().Err(err).Msg("terraform state doesn't exist, creating")
		utc := time.Now().UTC()
		aggregate, err = aggregates.NewTerraformState(uuid.MustParse(projectId), make([]byte, 0), utc, utc, make([]byte, 0))
		if err != nil {
			l.Warn().Err(err).Msg("cannot create state aggregate")
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
			return
		}
		err = repo.Add(aggregate)
		if err != nil {
			l.Warn().Err(err).Msg("Could not save to storage")
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save to storage"})
			return
		}
	}

	if err != nil {
		l.Warn().Err(err).Msg("Could not retrieve from storage")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return

	}

	storedLock, err := aggregate.GetLockInfo()
	if err != nil {
		l.Warn().Err(err)
		ctx.JSON(http.StatusLocked, gin.H{"message": "Already Locked"})
		return
	}
	if storedLock != nil {
		l.Warn().Msgf("Locked by different ID: %s", storedLock.ID)
		ctx.JSON(http.StatusLocked, gin.H{"message": "Already Locked", "ID": storedLock.ID})
		return
	}

	if err := aggregate.SetLockInfo(&reqLock); err != nil {
		l.Warn().Msg("Acquiring lock failed")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Acquiring lock failed"})
		return
	}

	layout := "2006-01-02T15:04:05.999999Z"
	createdTime, err := time.Parse(layout, reqLock.Created)

	if err != nil {
		l.Warn().Err(err).Msg("Error parsing timestamp")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Error parsing timestamp"})
		return
	}
	eventLogger := l.With().Str("eventType", string(models.Terraform)).Str("eventId", aggregate.GetID().String()).Str("type", "delayedMessage").Logger()
	c.delayTaskManger.Publish(
		events.LEASE_LOCK,
		&events.LeasedLock{
			Id:        aggregate.GetID(),
			Type:      models.Terraform,
			LeaseTime: createdTime,
		},
		time.Hour,
		context.Background(),
	)
	eventLogger.Info().Msg("Published delayed message")
	err = repo.Update(aggregate)
	if err != nil {
		l.Warn().Err(err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not save to storage"})
		return
	}

	ctx.JSON(http.StatusOK, reqLock.ID)
	return
}

// Unlock
// @BasePath /api/v1
// @Summary Unlock terraform state for a given project
// @Schemes http
// @Param projectId path string true "project id"
// @Description
// @Tags state
// @Accept json
// @Produce json
// @Success 200
func (c *StateController) Unlock(context *gin.Context) {
	projectId := context.Param("projectId")
	l := logger.GetWitContext(context).
		With().
		Str("projectId", projectId).
		Logger()

	var reqLock aggregates.LockInfo
	err := context.BindJSON(&reqLock)
	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Cannot bind object"})
		return
	}
	repo := context.MustGet(string(helpers.UnitOfWorkSetup)).(*repositories.UnitOfWork).TerraformStateDbRepository
	aggregate, err := repo.Get(uuid.MustParse(projectId))
	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": "Could not retrieve from storage"})
		return
	}

	err = aggregate.LeaseLock(&reqLock)

	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
	}
	err = repo.Update(aggregate)
	if err != nil {
		l.Warn().Err(err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
	}

	context.JSON(http.StatusOK, gin.H{"message": "ok"})
}
