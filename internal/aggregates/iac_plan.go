package aggregates

import (
	"github.com/google/uuid"
	"labraboard/internal/entities"
	"labraboard/internal/valueobjects/iacPlans"
)

type IaCPlanType string

var (
	Terraform IaCPlanType = "terraform"
	Tofu      IaCPlanType = "tofu"
)

type IacPlan struct {
	id            uuid.UUID
	changeSummary iacPlans.ChangeSummaryIacPlan
	changes       []iacPlans.ChangesIacPlan
	planType      IaCPlanType
	planJson      []byte
}

func NewIacPlan(id uuid.UUID, planType IaCPlanType, plan []byte) (*IacPlan, error) {
	return &IacPlan{
		id:       id,
		planType: planType,
		planJson: plan,
	}, nil
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
			plan.changeSummary = *summary

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
