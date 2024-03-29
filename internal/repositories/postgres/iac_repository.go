package postgres

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	vo "labraboard/internal/valueobjects"
)

type IaCRepository struct {
	database *Database
}

func NewIaCRepository(database *Database) (*IaCRepository, error) {
	return &IaCRepository{
		database: database,
	}, nil
}

func (repo *IaCRepository) Get(id uuid.UUID) (*aggregates.Iac, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get IaC")
	}

	return repo.Map(state)
}

func (repo *IaCRepository) Map(state *models.IaCDb) (*aggregates.Iac, error) {
	var envs []*vo.IaCEnv
	if state.Envs != nil {
		if err := json.Unmarshal(state.Envs, &envs); err != nil {
			return nil, errors.Wrap(err, "can't get envs on iac")
		}
	}
	var plans []*vo.Plans
	if state.Plans != nil {
		if err := json.Unmarshal(state.Plans, &plans); err != nil {
			return nil, errors.Wrap(err, "can't get plans on iac")
		}
	}
	var variables []*vo.IaCVariable
	if state.Variables != nil {
		if err := json.Unmarshal(state.Variables, &variables); err != nil {
			return nil, errors.Wrap(err, "can't get variables on iac")
		}
	}

	var iacRepo vo.IaCRepo
	if state.Repo != nil {
		if err := json.Unmarshal(state.Repo, &iacRepo); err != nil {
			return nil, errors.Wrap(err, "can't get repo on iac")
		}
	}

	iac, err := aggregates.NewIac(state.ID, vo.IaCType(state.IacType), plans, envs, iacRepo, variables)
	if err != nil {
		return nil, errors.Wrap(err, "can't create IaC Aggregate")
	}

	return iac, nil
}

func (repo *IaCRepository) Add(iac *aggregates.Iac) error {
	i, err := iac.Map()
	if err != nil {
		return errors.Wrap(err, "can't map IaC")
	}
	result := repo.database.GormDB.Create(i)
	return result.Error
}

func (repo *IaCRepository) Update(iac *aggregates.Iac) error {
	i, err := iac.Map()
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

func (repo *IaCRepository) GetAll() []*aggregates.Iac {

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
