package mappers

import (
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
)

type IacMapper[TDao *models.IaCDb, T *aggregates.Iac] struct {
}

func (i IacMapper[TDao, T]) Map(dao TDao) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (i IacMapper[TDao, T]) RevertMap(aggregate T) (TDao, error) {
	//TODO implement me
	panic("implement me")
}
