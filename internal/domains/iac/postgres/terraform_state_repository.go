package postgres

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/domains/iac/postgres/models"
)

type TerraformStateRepository struct {
	database Database
}

func NewTerraformStateRepository(database *Database) (*TerraformStateRepository, error) {
	return &TerraformStateRepository{
		database: *database,
	}, nil
}

func (repo *TerraformStateRepository) Get(id uuid.UUID) (*aggregates.TerraformState, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get state")
	}
	return aggregates.NewTerraformState(state.ProjectId, state.State, state.CreatedOn, state.ModifyOn, state.Lock)
}

func (repo *TerraformStateRepository) Add(state *aggregates.TerraformState) error {
	repo.database.GormDB.Create(&models.TerraformStateDb{
		ProjectId: state.GetID(),
		State:     state.GetByteState(),
		CreatedOn: state.CreatedOn,
		ModifyOn:  state.ModifyOn,
		Lock:      state.GetByteLock(),
	})
	return nil
}

func (repo *TerraformStateRepository) Update(state *aggregates.TerraformState) error {
	s, err := repo.getState(state.GetID())
	if err != nil {
		return errors.Wrap(err, "can't get state")
	}
	repo.database.GormDB.Model(&s).Updates(&models.TerraformStateDb{
		State:    state.GetByteState(),
		ModifyOn: state.ModifyOn,
		Lock:     state.GetByteLock(),
	})
	return nil
}

func (repo *TerraformStateRepository) getState(id uuid.UUID) (*models.TerraformStateDb, error) {
	var state models.TerraformStateDb
	repo.database.GormDB.First(&state, "project_id =?", id)
	return &state, nil
}
