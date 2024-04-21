package mappers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/postgres/models"
	"labraboard/internal/valueobjects/iacPlans"
)

type IacPlanMapper[TDao *models.IaCPlanDb, T *aggregates.IacPlan] struct {
}

func (i IacPlanMapper[TDao, T]) Map(dao *models.IaCPlanDb) (*aggregates.IacPlan, error) {
	var changes []iacPlans.ChangesIacPlan
	if dao.Changes != nil {
		if err := json.Unmarshal(dao.Changes, &changes); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal history config")
		}
	}
	var summary *iacPlans.ChangeSummaryIacPlan
	if dao.ChangeSummary != nil {
		if err := json.Unmarshal(dao.ChangeSummary, &summary); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal change summary")
		}
	}

	var historyConfig *iacPlans.HistoryProjectConfig
	if dao.Config != nil {
		if err := json.Unmarshal(dao.Config, &historyConfig); err != nil {
			return nil, errors.Wrap(err, "can't unmarshal history config")
		}
	}
	return aggregates.NewIacPlanExplicit(dao.ID, aggregates.IaCPlanType(dao.PlanType), historyConfig, summary, changes, dao.PlanJson)
}

func (i IacPlanMapper[TDao, T]) RevertMap(aggregate *aggregates.IacPlan) (*models.IaCPlanDb, error) {
	if aggregate == nil {
		return nil, errors.New("can't map nil IaC")
	}

	planJson, planType, planChanges, planChangeSummary := aggregate.Composite()

	changes, err := json.Marshal(planChanges)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall changes on planChanges")
	}

	summary, err := json.Marshal(planChangeSummary)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall changes on planChangeSummary")
	}

	config, err := json.Marshal(aggregate.HistoryConfig)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall config on historyConfig")
	}

	return &models.IaCPlanDb{
		ID:            aggregate.GetID(),
		ChangeSummary: summary,
		Changes:       changes,
		PlanJson:      planJson,
		PlanType:      string(planType),
		Config:        config,
	}, nil
}
