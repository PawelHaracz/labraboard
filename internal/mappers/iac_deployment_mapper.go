package mappers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	"labraboard/internal/valueobjects/iac"
)

type IacDeploymentMapper[TDao *models.IaCDeploymentDb, T *aggregates.IacDeployment] struct {
}

func (i IacDeploymentMapper[TDao, T]) Map(dao *models.IaCDeploymentDb) (*aggregates.IacDeployment, error) {
	var changes []iac.ChangesIac
	if dao.Changes != nil {
		if err := json.Unmarshal(dao.Changes, &changes); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal history config")
		}
	}
	var summary iac.ChangeSummaryIac
	if dao.ChangeSummary != nil {
		if err := json.Unmarshal(dao.ChangeSummary, &summary); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal change summary")
		}
	}

	var outputs []iac.Output
	if dao.Outputs != nil {
		if err := json.Unmarshal(dao.Outputs, &outputs); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal outputs")
		}
	}

	return aggregates.NewIacDeploymentExplicit(dao.ID, dao.PlanId, dao.ProjectId, dao.Started, dao.Deployed, dao.DeploymentType, changes, summary, outputs), nil
}

func (i IacDeploymentMapper[TDao, T]) RevertMap(aggregate *aggregates.IacDeployment) (*models.IaCDeploymentDb, error) {
	if aggregate == nil {
		return nil, errors.New("can't map nil IaC")
	}
	deploymentChanges, deploymentChangeSummary, deploymentOutputs := aggregate.Composite()
	deploymentId, planId, projectId, startedTime, deployedTime, deployedType := aggregate.GetMetadata()
	changes, err := json.Marshal(deploymentChanges)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall changes on changes")
	}

	summary, err := json.Marshal(deploymentChangeSummary)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall changes on changeSummary")
	}

	outputs, err := json.Marshal(deploymentOutputs)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall config on outputs")
	}

	return &models.IaCDeploymentDb{
		ID:             deploymentId,
		PlanId:         planId,
		ProjectId:      projectId,
		Started:        startedTime,
		Deployed:       deployedTime,
		DeploymentType: deployedType,
		ChangeSummary:  summary,
		Changes:        changes,
		Outputs:        outputs,
	}, nil
}
