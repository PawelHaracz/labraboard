package postgres

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/mappers"
	"labraboard/internal/repositories/postgres/models"
)

type IaCPlanRepository struct {
	database *Database
	mapper   mappers.Mapper[*models.IaCPlanDb, *aggregates.IacPlan]
}

func NewIaCPlanRepository(database *Database) (*IaCPlanRepository, error) {
	return &IaCPlanRepository{
		database: database,
		mapper:   mappers.IacPlanMapper[*models.IaCPlanDb, *aggregates.IacPlan]{},
	}, nil
}

func (repo *IaCPlanRepository) Get(id uuid.UUID) (*aggregates.IacPlan, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get IaC")
	}

	return repo.Map(state)
}

func (repo *IaCPlanRepository) Map(state *models.IaCPlanDb) (*aggregates.IacPlan, error) {

	iac, err := repo.mapper.Map(state)

	return iac, err
}

func (repo *IaCPlanRepository) Add(iac *aggregates.IacPlan) error {
	i, err := repo.mapper.RevertMap(iac)
	if err != nil {
		return errors.Wrap(err, "can't map IaC")
	}
	result := repo.database.GormDB.Create(i)
	return result.Error
}

func (repo *IaCPlanRepository) Update(iac *aggregates.IacPlan) error {
	i, err := repo.mapper.RevertMap(iac)
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

func (repo *IaCPlanRepository) GetAll() []*aggregates.IacPlan {
	var planDbs []*models.IaCPlanDb
	repo.database.GormDB.Find(&planDbs)
	plans := make([]*aggregates.IacPlan, len(planDbs))
	for _, plan := range planDbs {
		p, err := repo.Map(plan)
		if err != nil {
			//handle it
			continue
		}
		plans = append(plans, p)
	}
	return plans
}

func (repo *IaCPlanRepository) getState(id uuid.UUID) (*models.IaCPlanDb, error) {
	var state models.IaCPlanDb
	result := repo.database.GormDB.First(&state, "id =?", id)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "can't get state")
	}
	return &state, nil
}
