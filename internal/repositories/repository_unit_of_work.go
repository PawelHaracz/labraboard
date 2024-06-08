package repositories

import (
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/memory"
	db "labraboard/internal/repositories/postgres"
)

type UnitOfWorkConfiguration func(os *UnitOfWork) error

type UnitOfWork struct {
	TerraformStateDbRepository Repository[*aggregates.TerraformState]
	IacRepository              Repository[*aggregates.Iac]
	IacPlan                    Repository[*aggregates.IacPlan]
	IacDeployment              Repository[*aggregates.IacDeployment]
}

func NewUnitOfWork(configs ...UnitOfWorkConfiguration) (*UnitOfWork, error) {
	uow := &UnitOfWork{}
	for _, cfg := range configs {
		if err := cfg(uow); err != nil {
			return nil, err
		}
	}

	if uow.TerraformStateDbRepository == nil {
		return nil, errors.New("terraform state is not set")
	}
	if uow.IacRepository == nil {
		return nil, errors.New("iac Repository is not set")
	}
	if uow.IacPlan == nil {
		return nil, errors.New("iac plan Repository is not set")
	}
	if uow.IacDeployment == nil {
		return nil, errors.New("iac deployment Repository is not set")
	}
	return uow, nil
}

func WithTerraformStateDbRepository(database *db.Database) UnitOfWorkConfiguration {
	repository, err := db.NewTerraformStateRepository(database)
	if err != nil {
		return func(uow *UnitOfWork) error {
			return errors.Wrap(err, "can't create terraform state repository")
		}
	}

	return func(uow *UnitOfWork) error {
		uow.TerraformStateDbRepository = repository
		return nil
	}
}
func WithIaCRepositoryDbRepository(database *db.Database) UnitOfWorkConfiguration {
	repository, err := db.NewIaCRepository(database)
	if err != nil {
		return func(uow *UnitOfWork) error {
			return errors.Wrap(err, "can't create terraform state repository")
		}
	}

	return func(uow *UnitOfWork) error {
		uow.IacRepository = repository
		return nil
	}
}

func WithIacPlanRepositoryDbRepository(database *db.Database) UnitOfWorkConfiguration {
	repository, err := db.NewIaCPlanRepository(database)
	if err != nil {
		return func(uow *UnitOfWork) error {
			return errors.Wrap(err, "can't create terraform state repository")
		}
	}

	return func(uow *UnitOfWork) error {
		uow.IacPlan = repository
		return nil
	}
}

func WithIacDeploymentRepositoryDbRepository(database *db.Database) UnitOfWorkConfiguration {
	repository, err := db.NewIacDeployment(database)
	if err != nil {
		return func(uow *UnitOfWork) error {
			return errors.Wrap(err, "can't create terraform state repository")
		}
	}

	return func(uow *UnitOfWork) error {
		uow.IacDeployment = repository
		return nil
	}
}

func WithIacPlanRepositoryDbRepositoryMemory(repository interface{}) UnitOfWorkConfiguration {
	switch repo := repository.(type) {
	case *memory.GenericRepository[*aggregates.IacPlan]:
		return func(uow *UnitOfWork) error {
			uow.IacPlan = repo
			return nil
		}
	case *memory.GenericRepository[*aggregates.Iac]:
		return func(uow *UnitOfWork) error {
			uow.IacRepository = repo
			return nil
		}
	case *memory.GenericRepository[*aggregates.TerraformState]:
		return func(uow *UnitOfWork) error {
			uow.TerraformStateDbRepository = repo
			return nil
		}
	default:
		return func(uow *UnitOfWork) error {
			return errors.New("repository is not a IacPlanRepository")
		}
	}
}
