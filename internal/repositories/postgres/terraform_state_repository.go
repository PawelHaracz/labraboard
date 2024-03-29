package postgres

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
)

type TerraformStateRepository struct {
	database *Database
}

func NewTerraformStateRepository(database *Database) (*TerraformStateRepository, error) {
	return &TerraformStateRepository{
		database: database,
	}, nil
}

func (repo *TerraformStateRepository) Get(id uuid.UUID) (*aggregates.TerraformState, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get state")
	}
	return aggregates.NewTerraformState(state.ID, state.State, state.CreatedOn, state.ModifyOn, state.Lock)
}

func (repo *TerraformStateRepository) Add(state *aggregates.TerraformState) error {
	result := repo.database.GormDB.Create(&models.TerraformStateDb{
		ID:        state.GetID(),
		State:     state.GetByteState(),
		CreatedOn: state.CreatedOn,
		ModifyOn:  state.ModifyOn,
		Lock:      state.GetByteLock(),
	})
	return result.Error
}

func (repo *TerraformStateRepository) Update(state *aggregates.TerraformState) error {
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

func (repo *TerraformStateRepository) GetAll() []*aggregates.TerraformState {
	var dbs []*models.TerraformStateDb
	repo.database.GormDB.Find(&dbs)
	states := make([]*aggregates.TerraformState, len(dbs))
	for _, state := range dbs {
		p, err := aggregates.NewTerraformState(state.ID, state.State, state.CreatedOn, state.ModifyOn, state.Lock)
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
