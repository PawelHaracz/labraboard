package services

import (
	"labraboard/internal/aggregates"
	m "labraboard/internal/eventbus/memory"
	"labraboard/internal/logger"
	"labraboard/internal/managers"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/memory"
	"labraboard/internal/valueobjects"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

func TestNewIacService(t *testing.T) {
	ctx := context.Background()
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.Iac]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.IacPlan]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.TerraformState]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.IacDeployment]()))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	is, err := NewIacService(
		WithEventBus(m.NewMemoryEventBus()),
		WithUnitOfWork(uow),
		WithDelayTaskManagerPublisher(managers.NewMemoryDelayTask()))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if is == nil {
		t.Errorf("error: %v", "IacService is nil")
	}

	if is.publisher == nil {
		t.Errorf("error: %v", "IacService.planner is nil")
	}

	if is.unitOfWork == nil {
		t.Errorf("error: %v", "IacService.repositories is nil")
	}

	aggregate, _ := aggregates.NewIac(uuid.New(), valueobjects.Terraform, make([]*valueobjects.Plans, 0), make([]*valueobjects.IaCEnv, 0), nil, make([]*valueobjects.IaCVariable, 0))
	err = is.unitOfWork.IacRepository.Add(aggregate, ctx)

	if err != nil {
		t.Errorf("error during adding item: %v", err)
	}
}

func TestIacService_RunTerraformPlan(t *testing.T) {
	ctx := context.Background()
	logger.Init(1, true)
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.Iac]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.IacPlan]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.TerraformState]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.IacDeployment]()))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	is, err := NewIacService(
		WithEventBus(m.NewMemoryEventBus()),
		WithUnitOfWork(uow),
		WithDelayTaskManagerPublisher(managers.NewMemoryDelayTask()))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if is == nil {
		t.Errorf("error: %v", "IacService is nil")
	}

	if is.publisher == nil {
		t.Errorf("error: %v", "IacService.planner is nil")
	}

	if is.unitOfWork == nil {
		t.Errorf("error: %v", "IacService.repositories is nil")
	}

	aggregate, _ := aggregates.NewIac(uuid.New(), valueobjects.Terraform, make([]*valueobjects.Plans, 0), make([]*valueobjects.IaCEnv, 0), nil, make([]*valueobjects.IaCVariable, 0))
	err = is.unitOfWork.IacRepository.Add(aggregate, ctx)

	runner := TerraformPlanRunner{
		aggregate.GetID(), "", "", models.SHA, nil,
	}
	planId, err := is.RunTerraformPlan(runner, ctx)

	if planId == uuid.Nil {
		t.Errorf("error: %v details: %v", "planId is nil", err)
	}
}
