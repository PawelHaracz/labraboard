package mappers

import (
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	"labraboard/internal/valueobjects/iac"
	"testing"
	"time"
)

func TestIacDeploymentMapper_Map(t *testing.T) {
	//arrange
	var mapper = IacDeploymentMapper[*models.IaCDeploymentDb, *aggregates.IacDeployment]{}
	var db = &models.IaCDeploymentDb{
		ID:             uuid.New(),
		PlanId:         uuid.New(),
		ProjectId:      uuid.New(),
		Started:        time.Now(),
		Deployed:       time.Now().Add(5),
		DeploymentType: "terraform",
		ChangeSummary:  nil,
		Changes:        nil,
		Outputs:        nil,
	}

	act, err := mapper.Map(db)
	if err != nil {
		t.Fatal(err)
	}

	deploymentId, planId, projectId, startedTime, deployedTime, deployedType := act.GetMetadata()
	assert.Equal(t, deploymentId, db.ID)
	assert.Equal(t, planId, db.PlanId)
	assert.Equal(t, projectId, db.ProjectId)
	assert.Equal(t, startedTime, db.Started)
	assert.Equal(t, deployedTime, db.Deployed)
	assert.Equal(t, deployedType, db.DeploymentType)
}

func TestIacDeploymentMapper_RevertMap(t *testing.T) {
	//arrange
	var mapper = IacDeploymentMapper[*models.IaCDeploymentDb, *aggregates.IacDeployment]{}

	aggregate := aggregates.NewIacDeploymentExplicit(uuid.New(), uuid.New(), uuid.New(), time.Now().Add(10), time.Now().Add(20), "terraform", nil, iac.EmptySummary, nil)
	act, err := mapper.RevertMap(aggregate)
	if err != nil {
		t.Fatal(err)
	}

	deploymentId, planId, projectId, startedTime, deployedTime, deployedType := aggregate.GetMetadata()
	assert.Equal(t, deploymentId, act.ID)
	assert.Equal(t, planId, act.PlanId)
	assert.Equal(t, projectId, act.ProjectId)
	assert.Equal(t, startedTime, act.Started)
	assert.Equal(t, deployedTime, act.Deployed)
	assert.Equal(t, deployedType, act.DeploymentType)
}
