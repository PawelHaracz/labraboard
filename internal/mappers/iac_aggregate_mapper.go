package mappers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	vo "labraboard/internal/valueobjects"
)

type IacMapper[TDao *models.IaCDb, T *aggregates.Iac] struct {
}

func (i IacMapper[TDao, T]) Map(dao *models.IaCDb) (*aggregates.Iac, error) {
	var envs []*vo.IaCEnv
	if dao.Envs != nil {
		if err := json.Unmarshal(dao.Envs, &envs); err != nil {
			return nil, errors.Wrap(err, "can't get envs on iac")
		}
	}
	var plans []*vo.Plans
	if dao.Plans != nil {
		if err := json.Unmarshal(dao.Plans, &plans); err != nil {
			return nil, errors.Wrap(err, "can't get plans on iac")
		}
	}
	var variables []*vo.IaCVariable
	if dao.Variables != nil {
		if err := json.Unmarshal(dao.Variables, &variables); err != nil {
			return nil, errors.Wrap(err, "can't get variables on iac")
		}
	}

	var iacRepo *vo.IaCRepo
	if dao.Repo != nil {
		if err := json.Unmarshal(dao.Repo, &iacRepo); err != nil {
			return nil, errors.Wrap(err, "can't get repo on iac")
		}
	}

	iac, err := aggregates.NewIac(dao.ID, vo.IaCType(dao.IacType), plans, envs, iacRepo, variables)
	if err != nil {
		return nil, errors.Wrap(err, "can't create IaC Aggregate")
	}

	return iac, nil
}

func (i IacMapper[TDao, T]) RevertMap(aggregate *aggregates.Iac) (*models.IaCDb, error) {
	iacEnvs, iacVariables, repo := aggregate.Composite()
	iacRepo, err := json.Marshal(repo)

	if err != nil {
		return nil, errors.Wrap(err, "can't create repo on receiver")
	}

	envs, err := json.Marshal(iacEnvs)
	if err != nil {
		return nil, errors.Wrap(err, "can't create envs on receiver")
	}

	variables, err := json.Marshal(iacVariables)
	if err != nil {
		return nil, errors.Wrap(err, "can't create variables on receiver")
	}

	plans, err := json.Marshal(aggregate.GetPlans())
	if err != nil {
		return nil, errors.Wrap(err, "can't create plans on receiver")
	}

	return &models.IaCDb{
		ID:        aggregate.GetID(),
		IacType:   int(aggregate.IacType),
		Repo:      iacRepo,
		Envs:      envs,
		Variables: variables,
		Plans:     plans,
	}, nil
}
