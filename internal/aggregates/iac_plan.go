package aggregates

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/entities"
	"labraboard/internal/repositories/postgres/models"
	"labraboard/internal/valueobjects/iacPlans"
)

type IaCPlanType string

var (
	Terraform IaCPlanType = "terraform"
	Tofu      IaCPlanType = "tofu"
)

type IacPlan struct {
	id            uuid.UUID
	HistoryConfig *iacPlans.HistoryProjectConfig
	changeSummary *iacPlans.ChangeSummaryIacPlan
	changes       []iacPlans.ChangesIacPlan
	planType      IaCPlanType
	planJson      []byte
}

func NewIacPlan(id uuid.UUID, planType IaCPlanType, historyConfig *iacPlans.HistoryProjectConfig) (*IacPlan, error) {
	return &IacPlan{
		id:            id,
		planType:      planType,
		HistoryConfig: historyConfig,
	}, nil
}

func MapIacPlan(state *models.IaCPlanDb) *IacPlan {
	var summary *iacPlans.ChangeSummaryIacPlan
	if state.PlanJson != nil {
		if err := json.Unmarshal(state.PlanJson, &summary); err != nil {
			return nil
		}
	}
	var changes []iacPlans.ChangesIacPlan
	if state.Changes != nil {
		if err := json.Unmarshal(state.Changes, &changes); err != nil {
			return nil
		}
	}
	var historyConfig *iacPlans.HistoryProjectConfig
	if state.ChangeSummary != nil {
		if err := json.Unmarshal(state.ChangeSummary, &historyConfig); err != nil {
			return nil
		}
	}

	return &IacPlan{
		id:            state.ID,
		planType:      IaCPlanType(state.PlanType),
		HistoryConfig: historyConfig,
		changeSummary: summary,
		changes:       changes,
		planJson:      state.PlanJson,
	}
}

func (p *IacPlan) AddPlan(plan []byte) {
	p.planJson = plan
}

func newChangeIacPlanner(resourceType string, resourceName string, provider string, action iacPlans.PlanTypeAction) *iacPlans.ChangesIacPlan {
	return &iacPlans.ChangesIacPlan{
		ResourceType: resourceType,
		ResourceName: resourceName,
		Provider:     provider,
		Action:       action,
	}
}

func newChangeSummaryIacPlan(added int, changed int, removed int) *iacPlans.ChangeSummaryIacPlan {
	return &iacPlans.ChangeSummaryIacPlan{
		Add:    added,
		Change: changed,
		Remove: removed,
	}
}

// GetID returns the Iac root entity ID
func (plan *IacPlan) GetID() uuid.UUID {
	return plan.id
}

func (plan *IacPlan) AddChanges(plans ...entities.IacTerraformPlanJson) {
	var changes []iacPlans.ChangesIacPlan

	for _, p := range plans {
		if p.Type == entities.Version {
			continue
		}
		if p.Change == nil {
			if p.SummaryChanges == nil {
				continue
			}
			summary := newChangeSummaryIacPlan(p.SummaryChanges.Add, p.SummaryChanges.Change, p.SummaryChanges.Remove)
			plan.changeSummary = summary

		} else {
			planner := newChangeIacPlanner(p.Change.Resource.ResourceType, p.Change.Resource.ResourceName, p.Change.Resource.Provider, iacPlans.PlanTypeAction(p.Change.Action))
			changes = append(changes, *planner)
		}

	}

	plan.changes = changes
}

func (plan *IacPlan) GetChanges() (add int, change int, delete int) {
	return plan.changeSummary.Add, plan.changeSummary.Change, plan.changeSummary.Remove
}

func (plan *IacPlan) GetPlanJson() string {
	return string(plan.planJson)
}

func (plan *IacPlan) Map() (*models.IaCPlanDb, error) {
	if plan == nil {
		return nil, errors.New("can't map nil IaC")
	}
	changes, err := json.Marshal(plan.changes)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall changes on iac")
	}

	summary, err := json.Marshal(plan.changeSummary)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall changes on iac")
	}

	config, err := json.Marshal(plan.HistoryConfig)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshall config on iac")
	}
	return &models.IaCPlanDb{
		ID:            plan.id,
		ChangeSummary: summary,
		Changes:       changes,
		PlanJson:      plan.planJson,
		PlanType:      string(plan.planType),
		Config:        config,
	}, nil
}
