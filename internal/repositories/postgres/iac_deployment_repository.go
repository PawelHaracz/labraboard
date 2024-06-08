package postgres

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"labraboard/internal/aggregates"
	"labraboard/internal/mappers"
	"labraboard/internal/repositories/postgres/models"
)

type IacDeploymentRepository struct {
	database *Database
	mapper   mappers.Mapper[*models.IaCDeploymentDb, *aggregates.IacDeployment]
}

func NewIacDeployment(database *Database) (*IacDeploymentRepository, error) {
	return &IacDeploymentRepository{
		database: database,
		mapper:   mappers.IacDeploymentMapper[*models.IaCDeploymentDb, *aggregates.IacDeployment]{},
	}, nil
}

func (repo *IacDeploymentRepository) Get(id uuid.UUID, ctx context.Context) (*aggregates.IacDeployment, error) {
	state, err := repo.getState(id)
	if err != nil {
		return nil, errors.Wrap(err, "can't get IaC")
	}

	return repo.Map(state)
}

func (repo *IacDeploymentRepository) Map(state *models.IaCDeploymentDb) (*aggregates.IacDeployment, error) {

	iac, err := repo.mapper.Map(state)

	return iac, err
}

func (repo *IacDeploymentRepository) Add(iac *aggregates.IacDeployment, ctx context.Context) error {
	i, err := repo.mapper.RevertMap(iac)
	if err != nil {
		return errors.Wrap(err, "can't map IaC Plan")
	}
	result := repo.database.GormDB.Create(i)
	return result.Error
}

func (repo *IacDeploymentRepository) Update(iac *aggregates.IacDeployment, ctx context.Context) error {
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
	old.Deployed = i.Deployed
	old.Outputs = i.Outputs
	result := repo.database.GormDB.Save(&old)
	return result.Error
}

func (repo *IacDeploymentRepository) GetAll(ctx context.Context) []*aggregates.IacDeployment {
	var deploymentDbs []*models.IaCDeploymentDb
	repo.database.GormDB.Find(&deploymentDbs)
	plans := make([]*aggregates.IacDeployment, len(deploymentDbs))
	for _, deploymentDb := range deploymentDbs {
		p, err := repo.Map(deploymentDb)
		if err != nil {
			//handle it
			continue
		}
		plans = append(plans, p)
	}
	return plans
}

func (repo *IacDeploymentRepository) getState(id uuid.UUID) (*models.IaCDeploymentDb, error) {
	var state models.IaCDeploymentDb
	result := repo.database.GormDB.First(&state, "id =?", id)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "can't get state")
	}
	return &state, nil
}
