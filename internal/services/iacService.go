package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	"labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"labraboard/internal/managers"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
	"labraboard/internal/routers/api/dtos"
	vo "labraboard/internal/valueobjects"
	"time"
)

type TerraformPlanRunner struct {
	ProjectId  uuid.UUID
	Path       string
	Sha        string
	CommitType models.CommitType
	Variables  map[string]string
}

type IacConfiguration func(os *IacService) error

type IacService struct {
	publisher                 eventbus.EventPublisher
	unitOfWork                *repositories.UnitOfWork
	delayTaskManagerPublisher managers.DelayTaskManagerPublisher
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
		return nil, errors.New("unit of Work is not set")
	}
	if is.delayTaskManagerPublisher == nil {
		return nil, errors.New("delay task manager publisher is not set")
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

func WithDelayTaskManagerPublisher(delayTaskManagerPublisher managers.DelayTaskManagerPublisher) IacConfiguration {
	return func(is *IacService) error {
		is.delayTaskManagerPublisher = delayTaskManagerPublisher
		return nil
	}
}

func (svc *IacService) RunTerraformPlan(runner TerraformPlanRunner, ctx context.Context) (uuid.UUID, error) {
	planId := uuid.New()
	log := logger.GetWitContext(ctx).
		With().
		Str("planId", planId.String()).
		Str("gitSha", runner.Sha).
		Str("gitPath", runner.Path).
		Logger()

	iac, err := svc.unitOfWork.IacRepository.Get(runner.ProjectId, ctx)
	if err != nil {
		return uuid.Nil, err
	}

	log.Info().
		Msg("Creating project")

	iac.AddPlan(planId, runner.Sha, runner.Path, runner.Variables)
	err = svc.unitOfWork.IacRepository.Update(iac, ctx)
	if err != nil {
		log.Error().Err(err)
		return uuid.Nil, err
	}

	var event = events.PlanTriggered{
		ProjectId: runner.ProjectId,
		PlanId:    planId,
		RepoPath:  runner.Path,
		Commit: events.Commit{
			Type: runner.CommitType,
			Name: runner.Sha,
		},
		Variables:    runner.Variables,
		EnvVariables: iac.GetEnvs(false),
	}

	svc.publisher.Publish(events.TRIGGERED_PLAN, event, ctx)
	if err != nil {
		log.Error().Err(err)
		return uuid.Nil, err
	}

	return planId, nil
}

func (svc *IacService) GetProjects(ctx context.Context) ([]*aggregates.Iac, error) {
	return svc.unitOfWork.IacRepository.GetAll(ctx), nil // TODO implement pagination
}

func (svc *IacService) GetProject(projectId uuid.UUID, ctx context.Context) (*aggregates.Iac, error) {
	return svc.unitOfWork.IacRepository.Get(projectId, ctx)
}

func (svc *IacService) CreateProject(iacType vo.IaCType, repo *vo.IaCRepo, ctx context.Context) (uuid.UUID, error) {
	projectId := uuid.New()

	iac, err := aggregates.NewIac(projectId, iacType, make([]*vo.Plans, 0), make([]*vo.IaCEnv, 0), repo, make([]*vo.IaCVariable, 0))
	if err != nil {
		return uuid.Nil, err
	}

	if err = svc.unitOfWork.IacRepository.Add(iac, ctx); err != nil {
		return uuid.Nil, err
	}

	return projectId, nil
}

func (svc *IacService) GetPlans(projectId uuid.UUID, ctx context.Context) []*vo.Plans {
	iac, err := svc.unitOfWork.IacRepository.Get(projectId, ctx)
	if err != nil {
		return nil
	}
	return iac.GetPlans()
}

func (svc *IacService) GetPlan(projectId uuid.UUID, planId uuid.UUID, ctx context.Context) (*dtos.PlanWithOutputDto, error) {
	iac, err := svc.unitOfWork.IacRepository.Get(projectId, ctx)
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
		p, err := svc.unitOfWork.IacPlan.Get(plan.Id, ctx)
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

func (svc *IacService) AddEnv(projectId uuid.UUID, name string, value string, isSecret bool, ctx context.Context) error {
	iac, err := svc.GetProject(projectId, ctx)
	if err != nil {
		return err
	}

	if err = iac.AddEnv(name, value, isSecret); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac, ctx)
}

func (svc *IacService) RemoveEnv(projectId uuid.UUID, name string, ctx context.Context) error {
	iac, err := svc.GetProject(projectId, ctx)
	if err != nil {
		return err
	}

	if err = iac.RemoveEnv(name); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac, ctx)
}

func (svc *IacService) AddVariable(projectId uuid.UUID, name string, value string, ctx context.Context) error {
	iac, err := svc.GetProject(projectId, ctx)
	if err != nil {
		return err
	}

	if err = iac.SetVariable(name, value); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac, ctx)
}

func (svc *IacService) RemoveVariable(projectId uuid.UUID, name string, ctx context.Context) error {
	iac, err := svc.GetProject(projectId, ctx)
	if err != nil {
		return err
	}

	if err = iac.RemoveVariable(name); err != nil {
		return err
	}

	return svc.unitOfWork.IacRepository.Update(iac, ctx)
}

func (svc *IacService) SchedulePlan(projectId uuid.UUID, planId uuid.UUID, when time.Time, ctx context.Context) error {
	l := logger.GetWitContext(ctx).With().Time("when", when).Logger()
	if _, err := svc.unitOfWork.IacRepository.Get(projectId, ctx); err != nil {
		l.Error().Err(err)
		return errors.Wrap(err, fmt.Sprintf("Project doesn't exist: %s", projectId))
	}
	if _, err := svc.unitOfWork.IacPlan.Get(planId, ctx); err != nil {
		l.Warn().Err(err)
		return errors.Wrap(err, fmt.Sprintf("Plan doesn't exist: %s, please first run plan", planId))
	}
	now := time.Now()
	hasFuture := now.After(when)
	if !hasFuture {
		err := errors.New(fmt.Sprint("you cannot use time in past for schedule"))
		l.Warn().Err(err)
		return err
	}
	svc.delayTaskManagerPublisher.Publish(events.SCHEDULED_PLAN,
		events.ScheduledPlan{
			ProjectId: projectId,
			PlanId:    planId,
			When:      when,
		},
		when.Sub(now),
		l.WithContext(ctx))
	return nil
}

func (svc *IacService) ScheduleApply(projectId uuid.UUID, planId uuid.UUID, ctx context.Context) (uuid.UUID, error) {
	l := logger.GetWitContext(ctx).With().Logger()
	l.Info().Msg("Starting applying")

	if _, err := svc.unitOfWork.IacRepository.Get(projectId, ctx); err != nil {
		l.Error().Err(err)
		return uuid.Nil, errors.Wrap(err, fmt.Sprintf("Project doesn't exist: %s", projectId))
	}
	if _, err := svc.unitOfWork.IacPlan.Get(planId, ctx); err != nil {
		l.Warn().Err(err)
		return uuid.Nil, errors.Wrap(err, fmt.Sprintf("Plan doesn't exist: %s, please first run plan", planId))
	}

	var changeId = uuid.New()
	var event = &events.IacApplied{
		ChangeId:  changeId,
		ProjectId: projectId,
		PlanId:    planId,
		IacType:   aggregates.Terraform,
		Owner:     "Anonymous",
	}

	svc.publisher.Publish(events.IAC_APPLY_SCHEDULED, event, ctx)

	return changeId, nil
}
