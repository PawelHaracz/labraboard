package main

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	dbmemory "labraboard/internal/domains/iac/memory"
	"labraboard/internal/eventbus/events"
	vo "labraboard/internal/valueobjects"
	"testing"
)

func TestPlanTriggerHandler(t *testing.T) {
	db, _ := dbmemory.NewRepository()

	aggregate, _ := aggregates.NewIac(uuid.New(), vo.Terraform)
	if err := db.Add(aggregate); err != nil {
		t.Failed()
	}
	var planId = uuid.New()
	aggregate.AddPlan(planId)

	var obj = &events.PlanTriggered{
		ProjectId: aggregate.GetID(),
		PlanId:    planId,
	}
	handlePlanTriggered(db, *obj)
	//todo fix passing backend to plan

}
