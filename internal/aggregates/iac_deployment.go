package aggregates

import (
	"github.com/google/uuid"
	"labraboard/internal/entities"
	"labraboard/internal/valueobjects/iac"
	"time"
)

// todo design -> it should be combined somehow with IacPlan

type IacDeployment struct {
	id             uuid.UUID
	planId         uuid.UUID
	projectId      uuid.UUID
	deploymentType IaCDeploymentType
	startedTime    time.Time
	deployedTime   time.Time
	changes        []iac.ChangesIac
	changeSummary  iac.ChangeSummaryIac
	outputs        []iac.Output
}

func NewIacDeployment(id uuid.UUID, planId uuid.UUID, projectId uuid.UUID, deploymentType IaCDeploymentType) *IacDeployment {
	return &IacDeployment{
		id:             id,
		planId:         planId,
		projectId:      projectId,
		deploymentType: deploymentType,
		startedTime:    time.Now().UTC(),
		changes:        make([]iac.ChangesIac, 0),
		outputs:        make([]iac.Output, 0),
	}
}

func NewIacDeploymentExplicit(deploymentId uuid.UUID, planId uuid.UUID, projectId uuid.UUID, startedTime time.Time, deployedTime time.Time, deployedType string, changes []iac.ChangesIac, summary iac.ChangeSummaryIac, outputs []iac.Output) *IacDeployment {
	return &IacDeployment{
		id:             deploymentId,
		planId:         planId,
		projectId:      projectId,
		deploymentType: IaCDeploymentType(deployedType),
		startedTime:    startedTime,
		deployedTime:   deployedTime,
		changes:        changes,
		changeSummary:  summary,
		outputs:        outputs,
	}
}

// GetID returns the Iac root entity ID
func (deployment *IacDeployment) GetID() uuid.UUID {
	return deployment.id
}

func (deployment *IacDeployment) FinishDeployment(plans ...entities.IacTerraformOutputJson) {
	var changes []iac.ChangesIac
	deployment.deployedTime = time.Now().UTC()
	for _, p := range plans {
		if p.Type == entities.Version {
			continue
		}

		if p.Change == emptyPlanChange {

			if p.SummaryChanges == emptySummaryChange {
				continue
			}
			deployment.changeSummary = iac.ChangeSummaryIac{
				Add:    p.SummaryChanges.Add,
				Change: p.SummaryChanges.Change,
				Remove: p.SummaryChanges.Remove,
			}

		} else if p.Outputs != nil && len(p.Outputs) != 0 {
			for key, value := range p.Outputs {
				deployment.outputs = append(deployment.outputs, iac.Output{
					Name:      key,
					Sensitive: value.Sensitive,
					Type:      value.Type,
					Value:     value.Value,
				})
			}
		} else {
			planner := newChangeIacPlanner(p.Change.Resource.ResourceType, p.Change.Resource.ResourceName, p.Change.Resource.Provider, iac.PlanTypeAction(p.Change.Action))
			deployment.changes = append(changes, *planner)
		}
	}
}

func (deployment *IacDeployment) Composite() ([]iac.ChangesIac, iac.ChangeSummaryIac, []iac.Output) {
	return deployment.changes, deployment.changeSummary, deployment.outputs
}

func (deployment *IacDeployment) GetMetadata() (deploymentId uuid.UUID, planId uuid.UUID, projectId uuid.UUID, startedTime time.Time, deployedTime time.Time, deployedType string) {
	return deployment.id, deployment.planId, deployment.projectId, deployment.startedTime, deployment.deployedTime, string(deployment.deploymentType)
}
