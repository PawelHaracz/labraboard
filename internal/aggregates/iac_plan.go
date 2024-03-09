package aggregates

import (
	"github.com/google/uuid"
	"labraboard/internal/entities"
)

type PlanTypeAction string
type IaCPlanType string

var (
	Create PlanTypeAction = "create"
	Update PlanTypeAction = "update"
	Delete PlanTypeAction = "delete"
)

var (
	Terraform IaCPlanType = "terraform"
	Tofu      IaCPlanType = "tofu"
)

type IacPlan struct {
	id            uuid.UUID
	changeSummary ChangeSummaryIacPlan
	changes       []ChangeIacPlanner
	planType      IaCPlanType
}

type ChangeSummaryIacPlan struct {
	Add    int
	Change int
	Remove int
}

type ChangeIacPlanner struct {
	ResourceType string
	ResourceName string
	Provider     string
	Action       PlanTypeAction
}

func NewIacPlan(id uuid.UUID, planType IaCPlanType) (*IacPlan, error) {
	return &IacPlan{
		id:       id,
		planType: planType,
	}, nil
}

func newChangeIacPlanner(resourceType string, resourceName string, provider string, action PlanTypeAction) *ChangeIacPlanner {
	return &ChangeIacPlanner{
		resourceType,
		resourceName,
		provider,
		action,
	}
}

func newChangeSummaryIacPlan(added int, changed int, removed int) *ChangeSummaryIacPlan {
	return &ChangeSummaryIacPlan{
		added,
		changed,
		removed,
	}
}

func (plan *IacPlan) AddChanges(plans ...entities.IacTerraformPlanJson) {
	var changes []ChangeIacPlanner

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
			planner := newChangeIacPlanner(p.Change.Resource.ResourceType, p.Change.Resource.ResourceName, p.Change.Resource.Provider, PlanTypeAction(p.Change.Action))
			changes = append(changes, *planner)
		}

	}

	plan.changes = changes
}

func (plan *IacPlan) GetChanges() (add int, change int, delete int) {
	return plan.changeSummary.Add, plan.changeSummary.Change, plan.changeSummary.Remove
}
