package services

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	m "labraboard/internal/eventbus/memory"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/memory"
	"labraboard/internal/valueobjects"
	"testing"
)

func TestNewIacService(t *testing.T) {
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.Iac]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.IacPlan]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.TerraformState]()))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	is, err := NewIacService(
		WithEventBus(m.NewMemoryEventBus()),
		WithUnitOfWork(uow))

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
	err = is.unitOfWork.IacRepository.Add(aggregate)

	if err != nil {
		t.Errorf("error during adding item: %v", err)
	}
}

func TestIacService_RunTerraformPlan(t *testing.T) {
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.Iac]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.IacPlan]()),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(memory.NewGenericRepository[*aggregates.TerraformState]()))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	is, err := NewIacService(
		WithEventBus(m.NewMemoryEventBus()),
		WithUnitOfWork(uow))

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
	err = is.unitOfWork.IacRepository.Add(aggregate)

	planId, err := is.RunTerraformPlan(aggregate.GetID())

	if planId == uuid.Nil {
		t.Errorf("error: %v details: %v", "planId is nil", err)
	}
}
