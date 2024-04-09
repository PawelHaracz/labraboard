package mappers

import (
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
)

type TerraformStatenMapper[TDao *models.TerraformStateDb, T *aggregates.TerraformState] struct {
}

func (i TerraformStatenMapper[TDao, T]) Map(dao TDao) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (i TerraformStatenMapper[TDao, T]) RevertMap(aggregate T) (TDao, error) {
	//TODO implement me
	panic("implement me")
}
