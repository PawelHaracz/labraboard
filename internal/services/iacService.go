package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/repositories"
	"labraboard/internal/routers/api/dtos"
	vo "labraboard/internal/valueobjects"
)

type IacConfiguration func(os *IacService) error

type IacService struct {
	publisher         eventbus.EventPublisher
	repository        repositories.Repository[*aggregates.Iac]
	iacPlanRepository repositories.Repository[*aggregates.IacPlan]
}

func NewIacService(configs ...IacConfiguration) (*IacService, error) {
	is := &IacService{}
	for _, cfg := range configs {
		if err := cfg(is); err != nil {
			return nil, err
		}
	}

	if is.publisher == nil {
		return nil, errors.New("publisher is not set")
	}
	if is.repository == nil {
		return nil, errors.New("repositories is not set")
	}

	return is, nil
}

func WithEventBus(eb eventbus.EventPublisher) IacConfiguration {
	return func(is *IacService) error {
		is.publisher = eb
		return nil
	}
}

func WithRepository(r repositories.Repository[*aggregates.Iac]) IacConfiguration {
	return func(is *IacService) error {
		is.repository = r
		return nil
	}
}

func (svc *IacService) RunTerraformPlan(projectId uuid.UUID) (uuid.UUID, error) {
	planId := uuid.New()

	iac, err := svc.repository.Get(projectId)
	if err != nil {
		return uuid.Nil, err
	}

	iac.AddPlan(planId)
	err = svc.repository.Update(iac)
	if err != nil {
		return uuid.Nil, err
	}

	var event = events.PlanTriggered{
		ProjectId: projectId,
		PlanId:    planId}

	svc.publisher.Publish(eventbus.TRIGGERED_PLAN, event, context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	return planId, nil
}

func (svc *IacService) GetProjects() ([]*aggregates.Iac, error) {
	return svc.repository.GetAll(), nil // TODO implement pagination
}

func (svc *IacService) GetProject(projectId uuid.UUID) (*aggregates.Iac, error) {
	return svc.repository.Get(projectId)
}

func (svc *IacService) CreateProject(iacType vo.IaCType, repo *vo.IaCRepo) (uuid.UUID, error) {
	projectId := uuid.New()

	iac, err := aggregates.NewIac(projectId, iacType, make([]*vo.Plans, 0), make([]*vo.IaCEnv, 0), repo, make([]*vo.IaCVariable, 0))
	if err != nil {
		return uuid.Nil, err
	}

	if err := svc.repository.Add(iac); err != nil {
		return uuid.Nil, err
	}

	return projectId, nil
}

func (svc *IacService) GetPlans(projectId uuid.UUID) []*vo.Plans {
	iac, err := svc.repository.Get(projectId)
	if err != nil {
		return nil
	}
	return iac.GetPlans()
}

func (svc *IacService) GetPlan(projectId uuid.UUID, planId uuid.UUID) (*dtos.PlanWithOutputDto, error) {
	iac, err := svc.repository.Get(projectId)
	if err != nil {
		return nil, err
	}

	plan, err := iac.GetPlan(planId)
	if err != nil {
		return nil, err
	}

	result := &dtos.PlanWithOutputDto{
		Id:        plan.Id.String(),
		CreatedOn: plan.CreatedOn,
		Status:    string(plan.Status),
	}

	p, err := svc.iacPlanRepository.Get(plan.Id)
	if err == nil {
		m := map[string]interface{}{
			"changes": p.GetChanges(),
			"json":    p.GetPlanJson(),
		}
		result.Outputs = m

		p.GetPlanJson()
	}
	return result, nil
}
