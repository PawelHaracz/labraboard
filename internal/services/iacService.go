package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
	"labraboard/internal/routers/api/dtos"
	vo "labraboard/internal/valueobjects"
)

type IacConfiguration func(os *IacService) error

type IacService struct {
	publisher  eventbus.EventPublisher
	unitOfWork *repositories.UnitOfWork
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
	if is.unitOfWork == nil {
		return nil, errors.New("Unit of Work is not set")
	}

	return is, nil
}

func WithEventBus(eb eventbus.EventPublisher) IacConfiguration {
	return func(is *IacService) error {
		is.publisher = eb
		return nil
	}
}

func WithUnitOfWork(r *repositories.UnitOfWork) IacConfiguration {
	return func(is *IacService) error {
		is.unitOfWork = r
		return nil
	}
}

func (svc *IacService) RunTerraformPlan(projectId uuid.UUID, path string, sha string, commitType models.CommitType, variables map[string]string) (uuid.UUID, error) {
	planId := uuid.New()

	iac, err := svc.unitOfWork.IacRepository.Get(projectId)
	if err != nil {
		return uuid.Nil, err
	}

	iac.AddPlan(planId, sha, path, variables)
	err = svc.unitOfWork.IacRepository.Update(iac)
	if err != nil {
		return uuid.Nil, err
	}

	var event = events.PlanTriggered{
		ProjectId: projectId,
		PlanId:    planId,
		RepoPath:  path,
		Commit: events.Commit{
			Type: commitType,
			Name: sha,
		},
		Variables: variables,
	}

	svc.publisher.Publish(events.TRIGGERED_PLAN, event, context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	return planId, nil
}

func (svc *IacService) GetProjects() ([]*aggregates.Iac, error) {
	return svc.unitOfWork.IacRepository.GetAll(), nil // TODO implement pagination
}

func (svc *IacService) GetProject(projectId uuid.UUID) (*aggregates.Iac, error) {
	return svc.unitOfWork.IacRepository.Get(projectId)
}

func (svc *IacService) CreateProject(iacType vo.IaCType, repo *vo.IaCRepo) (uuid.UUID, error) {
	projectId := uuid.New()

	iac, err := aggregates.NewIac(projectId, iacType, make([]*vo.Plans, 0), make([]*vo.IaCEnv, 0), repo, make([]*vo.IaCVariable, 0))
	if err != nil {
		return uuid.Nil, err
	}

	if err := svc.unitOfWork.IacRepository.Add(iac); err != nil {
		return uuid.Nil, err
	}

	return projectId, nil
}

func (svc *IacService) GetPlans(projectId uuid.UUID) []*vo.Plans {
	iac, err := svc.unitOfWork.IacRepository.Get(projectId)
	if err != nil {
		return nil
	}
	return iac.GetPlans()
}

func (svc *IacService) GetPlan(projectId uuid.UUID, planId uuid.UUID) (*dtos.PlanWithOutputDto, error) {
	iac, err := svc.unitOfWork.IacRepository.Get(projectId)
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
	if plan.Status == vo.Succeed {
		p, err := svc.unitOfWork.IacPlan.Get(plan.Id)
		add, update, deleteItem := p.GetChanges()
		if err == nil {
			m := map[string]interface{}{
				"changes": map[string]interface{}{
					"add":    add,
					"update": update,
					"delete": deleteItem,
				},
				"json": p.GetPlanJson(),
			}
			result.Outputs = m

			p.GetPlanJson()
		}
	}
	return result, nil
}

func (svc *IacService) AddEnv(projectId uuid.UUID, name string, value string, isSecret bool) error {
	iac, err := svc.GetProject(projectId)
	if err != nil {
		return err
	}

	if err = iac.AddEnv(name, value, isSecret); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac)
}

func (svc *IacService) RemoveEnv(projectId uuid.UUID, name string) error {
	iac, err := svc.GetProject(projectId)
	if err != nil {
		return err
	}

	if err = iac.RemoveEnv(name); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac)
}

func (svc *IacService) AddVariable(projectId uuid.UUID, name string, value string) error {
	iac, err := svc.GetProject(projectId)
	if err != nil {
		return err
	}

	if err = iac.SetVariable(name, value); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac)
}

func (svc *IacService) RemoveVariable(projectId uuid.UUID, name string) error {
	iac, err := svc.GetProject(projectId)
	if err != nil {
		return err
	}

	if err = iac.RemoveVariable(name); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac)
}
