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

var (
	emptyPlanChange    = entities.IacTerraformPlanChangeJson{}
	emptySummaryChange = entities.IacTerraformPlanSummaryChangesJson{}
)

type IacPlan struct {
	id            uuid.UUID
	HistoryConfig *iacPlans.HistoryProjectConfig
	changeSummary *iacPlans.ChangeSummaryIacPlan
	changes       []iacPlans.ChangesIacPlan
	planType      IaCPlanType
	planJson      []byte
	planRaw       []byte
}

func NewIacPlan(id uuid.UUID, planType IaCPlanType, historyConfig *iacPlans.HistoryProjectConfig) (*IacPlan, error) {
	return &IacPlan{
		id:            id,
		planType:      planType,
		HistoryConfig: historyConfig,
	}, nil
}

func NewIacPlanExplicit(id uuid.UUID, planType IaCPlanType, config *iacPlans.HistoryProjectConfig, summary *iacPlans.ChangeSummaryIacPlan, changes []iacPlans.ChangesIacPlan, planJson []byte, planRaw []byte) (*IacPlan, error) {
	return &IacPlan{
		id:            id,
		planType:      planType,
		HistoryConfig: config,
		changeSummary: summary,
		changes:       changes,
		planJson:      planJson,
		planRaw:       planRaw,
	}, nil
}

func (p *IacPlan) AddPlan(plan []byte, planRaw []byte) {
	p.planJson = plan
	p.planRaw = planRaw
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

		if p.Change == emptyPlanChange {

			if p.SummaryChanges == emptySummaryChange {
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

func (plan *IacPlan) Composite() (planJson []byte, planType IaCPlanType, changes []iacPlans.ChangesIacPlan, summary iacPlans.ChangeSummaryIacPlan, planRaw []byte) {
	if plan.changeSummary == nil {
		return plan.planJson, plan.planType, plan.changes, iacPlans.ChangeSummaryIacPlan{Add: 0, Change: 0, Remove: 0}, plan.planRaw
	}
	return plan.planJson, plan.planType, plan.changes, *plan.changeSummary, plan.planRaw
}

func (p *IacPlan) GetPlanType() string {
	return string(p.planType)
}

func (p *IacPlan) GetPlanRaw() []byte { return p.planRaw }
