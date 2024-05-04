package postgres

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"labraboard/internal/aggregates"
	"labraboard/internal/mappers"
	"labraboard/internal/repositories/postgres/models"
)

type TerraformStateRepository struct {
	database *Database
	mapper   mappers.Mapper[*models.TerraformStateDb, *aggregates.TerraformState]
}

func NewTerraformStateRepository(database *Database) (*TerraformStateRepository, error) {
	return &TerraformStateRepository{
		database: database,
		mapper:   mappers.TerraformStatenMapper[*models.TerraformStateDb, *aggregates.TerraformState]{},
	}, nil
}

func (repo *TerraformStateRepository) Get(id uuid.UUID, ctx context.Context) (*aggregates.TerraformState, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get state")
	}
	return repo.mapper.Map(state)
	//return aggregates.NewTerraformState(state.ID, state.State, state.CreatedOn, state.ModifyOn, state.Lock)
}

func (repo *TerraformStateRepository) Add(state *aggregates.TerraformState, ctx context.Context) error {
	model, err := repo.mapper.RevertMap(state)
	if err != nil {
		return errors.Wrap(err, "can't map state")
	}
	result := repo.database.GormDB.Create(model)
	return result.Error
}

func (repo *TerraformStateRepository) Update(state *aggregates.TerraformState, ctx context.Context) error {
	s, err := repo.getState(state.GetID())
	if err != nil {
		return errors.Wrap(err, "can't get state")
	}

	s.State = state.GetByteState()
	s.ModifyOn = state.ModifyOn
	s.Lock = state.GetByteLock()
	result := repo.database.GormDB.Save(&s)
	return result.Error
}

func (repo *TerraformStateRepository) GetAll(ctx context.Context) []*aggregates.TerraformState {
	var dbs []*models.TerraformStateDb
	repo.database.GormDB.Find(&dbs)
	states := make([]*aggregates.TerraformState, len(dbs))
	for _, state := range dbs {
		p, err := repo.mapper.Map(state)
		if err != nil {
			//handle it
			continue
		}
		states = append(states, p)
	}
	return states
}

func (repo *TerraformStateRepository) getState(id uuid.UUID) (*models.TerraformStateDb, error) {
	var state models.TerraformStateDb
	result := repo.database.GormDB.First(&state, "id =?", id)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "can't get state")
	}
	return &state, nil
}
