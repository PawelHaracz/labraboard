package mappers

import (
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	vo "labraboard/internal/valueobjects"
	"testing"
)

var mapper = IacMapper[*models.IaCDb, *aggregates.Iac]{}

func TestIacMapper_Map(t *testing.T) {
	//arrange
	var db = &models.IaCDb{
		ID:        uuid.New(),
		IacType:   0,
		Repo:      nil,
		Envs:      nil,
		Variables: nil,
		Plans:     nil,
	}

	act, err := mapper.Map(db)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, act.GetID(), db.ID)
	assert.Equal(t, int(act.IacType), db.IacType)
}

func TestIacMapper_RevertMap(t *testing.T) {
	var aggregate, _ = aggregates.NewIac(uuid.New(), vo.IaCType(0), nil, nil, nil, nil)
	act, err := mapper.RevertMap(aggregate)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, aggregate.GetID(), act.ID)
	assert.Equal(t, int(aggregate.IacType), act.IacType)
}
