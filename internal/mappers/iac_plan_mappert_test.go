package mappers

import (
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	"testing"
)

func TestIacPlanMapper_Map(t *testing.T) {
	//arrange
	var mapper = IacPlanMapper[*models.IaCPlanDb, *aggregates.IacPlan]{}
	var db = &models.IaCPlanDb{
		ID:            uuid.New(),
		ChangeSummary: nil,
		PlanType:      "terraform",
		Config:        nil,
		PlanJson:      nil,
		Changes:       nil,
	}

	act, err := mapper.Map(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, act.GetID(), db.ID)
	assert.Equal(t, act.GetPlanType(), db.PlanType)
}

func TestIacPlanMapper_RevertMap(t *testing.T) {
	//arrange
	var mapper = IacPlanMapper[*models.IaCPlanDb, *aggregates.IacPlan]{}

	aggregate, _ := aggregates.NewIacPlan(uuid.New(), "terraform", nil)

	act, err := mapper.RevertMap(aggregate)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, aggregate.GetID(), act.ID)
	assert.Equal(t, aggregate.GetPlanType(), act.PlanType)
}
