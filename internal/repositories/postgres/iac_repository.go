package postgres

import (
	_ "encoding/json"
	"labraboard/internal/aggregates"
	"labraboard/internal/mappers"
	"labraboard/internal/repositories/postgres/models"
	_ "labraboard/internal/valueobjects"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type IaCRepository struct {
	database *Database
	mapper   mappers.Mapper[*models.IaCDb, *aggregates.Iac]
}

func NewIaCRepository(database *Database) (*IaCRepository, error) {
	return &IaCRepository{
		database: database,
		mapper:   mappers.IacMapper[*models.IaCDb, *aggregates.Iac]{},
	}, nil
}

func (repo *IaCRepository) Get(id uuid.UUID, ctx context.Context) (*aggregates.Iac, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get IaC")
	}

	return repo.Map(state)
}

func (repo *IaCRepository) Map(state *models.IaCDb) (*aggregates.Iac, error) {
	return repo.mapper.Map(state)
}

func (repo *IaCRepository) Add(iac *aggregates.Iac, ctx context.Context) error {
	i, err := repo.mapper.RevertMap(iac)
	if err != nil {
		return errors.Wrap(err, "can't map IaC")
	}
	result := repo.database.GormDB.Create(i)
	return result.Error
}

func (repo *IaCRepository) Update(iac *aggregates.Iac, ctx context.Context) error {
	i, err := repo.mapper.RevertMap(iac)
	if err != nil {
		return errors.Wrap(err, "can't map IaC")
	}
	old, err := repo.getState(iac.GetID())
	if err != nil {
		return errors.Wrap(err, "can't get state")
	}

	old.Repo = i.Repo
	old.Envs = i.Envs
	old.Variables = i.Variables
	old.Plans = i.Plans
	result := repo.database.GormDB.Save(&old)
	return result.Error
}

func (repo *IaCRepository) GetAll(ctx context.Context) []*aggregates.Iac {

	var dbs []*models.IaCDb
	repo.database.GormDB.Find(&dbs)
	iacs := make([]*aggregates.Iac, len(dbs))
	for _, db := range dbs {
		p, err := repo.Map(db)
		if err != nil {
			//handle it
			continue
		}
		iacs = append(iacs, p)
	}
	return iacs

}

func (repo *IaCRepository) getState(id uuid.UUID) (*models.IaCDb, error) {
	var state models.IaCDb
	result := repo.database.GormDB.First(&state, "id =?", id)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "can't get state")
	}
	return &state, nil
}
