package mappers

import (
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
)

type TerraformStatenMapper[TDao *models.TerraformStateDb, T *aggregates.TerraformState] struct {
}

func (i TerraformStatenMapper[TDao, T]) Map(dao *models.TerraformStateDb) (*aggregates.TerraformState, error) {
	return aggregates.NewTerraformState(dao.ID, dao.State, dao.CreatedOn, dao.ModifyOn, dao.Lock)
}

func (i TerraformStatenMapper[TDao, T]) RevertMap(aggregate *aggregates.TerraformState) (*models.TerraformStateDb, error) {
	return &models.TerraformStateDb{
		ID:        aggregate.GetID(),
		State:     aggregate.GetByteState(),
		CreatedOn: aggregate.CreatedOn,
		ModifyOn:  aggregate.ModifyOn,
		Lock:      aggregate.GetByteLock(),
	}, nil
}
