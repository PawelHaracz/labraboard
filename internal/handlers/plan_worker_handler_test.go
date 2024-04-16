package handlers

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/repositories"
	dbmemory "labraboard/internal/repositories/memory"
	vo "labraboard/internal/valueobjects"
	"testing"
)

func TestPlanTriggerHandler(t *testing.T) {
	uow, _ := repositories.NewUnitOfWork(
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.IacPlan](),
		),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.Iac](),
		),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.TerraformState](),
		),
	)

	aggregate, _ := aggregates.NewIac(uuid.New(), vo.Terraform, nil, nil, nil, nil)
	if err := uow.IacRepository.Add(aggregate); err != nil {
		t.Failed()
	}
	var planId = uuid.New()
	aggregate.AddPlan(planId)
	aggregate.AddEnv("ARM_TENANT_ID", "4c83ec3e-26b4-444f-afb7-8b171cd1b420", false)
	aggregate.AddEnv("ARM_CLIENT_ID", "99cc9476-40fd-48b6-813f-e79e0ff830fc", false)
	aggregate.AddEnv("ARM_CLIENT_SECRET", "TODO", true)
	aggregate.AddEnv("ARM_SUBSCRIPTION_ID", "cb5863b1-784d-4813-b2c7-e87919081ecb", false)

	aggregate.AddRepo("https://github.com/microsoft/terraform-azure-devops-starter.git", "master", "101-terraform-job/terraform")

	aggregate.SetVariable("environment", "staging")
	aggregate.SetVariable("location", "Poland Center")

	var obj = &events.PlanTriggered{
		ProjectId: aggregate.GetID(),
		PlanId:    planId,
	}
	handlePlanTriggered(uow, *obj)
	aggregate, _ = uow.IacRepository.Get(aggregate.GetID())
	plan, err := aggregate.GetPlan(planId)
	if err != nil {
		t.Errorf("can't fetch plan: %v", err)
	}

	if plan.Status != vo.Succeed {
		t.Errorf("Plan Status not set to Succeed")
	}

	planAggregate, _ := uow.IacPlan.Get(planId)
	planAggregate.GetChanges()
	planJson := planAggregate.GetPlanJson()
	if planJson == "" {
		t.Errorf("Plan Json not set")
	}
}
