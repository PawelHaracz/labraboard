package mappers

import "labraboard/internal/aggregates"

type Mapper[TDao any, T aggregates.Aggregate] interface {
	Map(dao TDao) (T, error)
	RevertMap(aggregate T) (TDao, error)
}
