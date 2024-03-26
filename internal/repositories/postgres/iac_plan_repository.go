package postgres

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	"labraboard/internal/valueobjects/iacPlans"
)

type IaCPlanRepository struct {
	database *Database
}

func NewIaCPlanRepository(database *Database) (*IaCPlanRepository, error) {
	return &IaCPlanRepository{
		database: database,
	}, nil
}

func (repo *IaCPlanRepository) Get(id uuid.UUID) (*aggregates.IacPlan, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get IaC")
	}

	var summary iacPlans.ChangeSummaryIacPlan
	if state.PlanJson != nil {
		if err := json.Unmarshal(state.PlanJson, &summary); err != nil {
			return nil, errors.Wrap(err, "can't get envs on iac")
		}
	}
	var changes []iacPlans.ChangesIacPlan
	if state.Changes != nil {
		if err := json.Unmarshal(state.Changes, &changes); err != nil {
			return nil, errors.Wrap(err, "can't get plans on iac")
		}
	}

	iac, err := aggregates.NewIacPlan(state.ID, aggregates.IaCPlanType(state.PlanType), state.PlanJson, &summary, changes)
	if err != nil {
		return nil, errors.Wrap(err, "can't create IaC Aggregate")
	}

	return iac, nil
}

func (repo *IaCPlanRepository) Add(iac *aggregates.IacPlan) error {
	i, err := iac.Map()
	if err != nil {
		return errors.Wrap(err, "can't map IaC")
	}
	result := repo.database.GormDB.Create(i)
	return result.Error
}

func (repo *IaCPlanRepository) Update(iac *aggregates.IacPlan) error {
	i, err := iac.Map()
	if err != nil {
		return errors.Wrap(err, "can't map IaC")
	}
	old, err := repo.getState(iac.GetID())
	if err != nil {
		return errors.Wrap(err, "can't get state")
	}

	old.Changes = i.Changes
	old.ChangeSummary = i.ChangeSummary
	old.PlanJson = i.PlanJson
	result := repo.database.GormDB.Save(&old)
	return result.Error
}

func (repo *IaCPlanRepository) getState(id uuid.UUID) (*models.IaCPlanDb, error) {
	var state models.IaCPlanDb
	result := repo.database.GormDB.First(&state, "id =?", id)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "can't get state")
	}
	return &state, nil
}
