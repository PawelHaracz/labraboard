package mappers

import (
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
)

type IacPlanMapper[TDao *models.IaCPlanDb, T *aggregates.IacPlan] struct {
}

func (i IacPlanMapper[TDao, T]) Map(dao TDao) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (i IacPlanMapper[TDao, T]) RevertMap(aggregate T) (TDao, error) {
	//TODO implement me
	panic("implement me")
}
