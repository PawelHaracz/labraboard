package main

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/eventbus/events"
	dbmemory "labraboard/internal/repositories/memory"
	vo "labraboard/internal/valueobjects"
	"testing"
)

func TestPlanTriggerHandler(t *testing.T) {
	db, _ := dbmemory.NewRepository()

	aggregate, _ := aggregates.NewIac(uuid.New(), vo.Terraform, nil, nil, vo.IaCRepo{}, nil)
	if err := db.Add(aggregate); err != nil {
		t.Failed()
	}
	var planId = uuid.New()
	aggregate.AddPlan(planId)
	aggregate.AddEnv("ARM_TENANT_ID", "4c83ec3e-26b4-444f-afb7-8b171cd1b420", false)
	aggregate.AddEnv("ARM_CLIENT_ID", "99cc9476-40fd-48b6-813f-e79e0ff830fc", false)
	aggregate.AddEnv("ARM_CLIENT_SECRET", "CeP8Q~yoYHlWeEw_WkgmH85rHT6ur.7s_UY9JclB", true)
	aggregate.AddEnv("ARM_SUBSCRIPTION_ID", "cb5863b1-784d-4813-b2c7-e87919081ecb", false)

	aggregate.AddRepo("https://github.com/microsoft/terraform-azure-devops-starter.git", "master", "101-terraform-job/terraform")

	aggregate.SetVariable("environment", "staging")
	aggregate.SetVariable("location", "Poland Center")

	var obj = &events.PlanTriggered{
		ProjectId: aggregate.GetID(),
		PlanId:    planId,
	}
	handlePlanTriggered(db, *obj)
	aggregate, _ = db.Get(aggregate.GetID())
	plan, err := aggregate.GetPlan(planId)
	if err != nil {
		t.Errorf("can't fetch plan: %v", err)
	}

	if plan.Status != vo.Succeed {
		t.Errorf("Plan Status not set to Succeed")
	}

	planAggregate, _ := db.GetPlan(planId)
	planAggregate.GetChanges()
	planJson := planAggregate.GetPlanJson()
	if planJson == "" {
		t.Errorf("Plan Json not set")
	}

}
