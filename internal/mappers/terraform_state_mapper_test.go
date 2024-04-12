package mappers

import (
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	"testing"
	"time"
)

func TestTerraformStatenMapper_Map(t *testing.T) {
	//arrange
	var mapper = TerraformStatenMapper[*models.TerraformStateDb, *aggregates.TerraformState]{}

	db := &models.TerraformStateDb{
		ID:        uuid.New(),
		State:     make([]byte, 0),
		CreatedOn: time.Now(),
		ModifyOn:  time.Now(),
		Lock:      make([]byte, 0),
	}
	db.Lock = append(db.Lock, 1)
	db.State = append(db.State, 2)

	//act
	act, err := mapper.Map(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, act.GetID(), db.ID)
	assert.Equal(t, act.CreatedOn, db.CreatedOn)
	assert.Equal(t, act.ModifyOn, db.ModifyOn)
}

func TestTerraformStatenMapper_RevertMap(t *testing.T) {
	var mapper = TerraformStatenMapper[*models.TerraformStateDb, *aggregates.TerraformState]{}
	aggregate, err := aggregates.NewTerraformState(uuid.New(), make([]byte, 0), time.Now(), time.Now(), make([]byte, 0))

	if err != nil {
		t.Fatal(err)
	}
	act, err := mapper.RevertMap(aggregate)

	assert.Equal(t, aggregate.GetID(), act.ID)
	assert.Equal(t, aggregate.CreatedOn, act.CreatedOn)
	assert.Equal(t, aggregate.ModifyOn, act.ModifyOn)
}
