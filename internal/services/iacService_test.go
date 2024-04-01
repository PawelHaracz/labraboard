package services

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/repositories/memory"
	"labraboard/internal/valueobjects"
	"testing"
)

func TestNewIacService(t *testing.T) {
	//tfPlanner, _ := models.NewTerraformPlanner()
	r := memory.NewGenericRepository[*aggregates.Iac]()

	is, err := NewIacService(
		//WithPlanner(tfPlanner),
		WithRepository(r))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if is == nil {
		t.Errorf("error: %v", "IacService is nil")
	}

	if is.publisher == nil {
		t.Errorf("error: %v", "IacService.planner is nil")
	}

	if is.repository == nil {
		t.Errorf("error: %v", "IacService.repositories is nil")
	}

	aggregate, _ := aggregates.NewIac(uuid.New(), valueobjects.Terraform, make([]*valueobjects.Plans, 0), make([]*valueobjects.IaCEnv, 0), nil, make([]*valueobjects.IaCVariable, 0))
	err = is.repository.Add(aggregate)

	if err != nil {
		t.Errorf("error during adding item: %v", err)
	}
}

func TestIacService_RunTerraformPlan(t *testing.T) {
	//tfPlanner, _ := models.NewTerraformPlanner()
	r := memory.NewGenericRepository[*aggregates.Iac]()

	is, _ := NewIacService(
		//WithPlanner(tfPlanner),
		WithRepository(r))

	projectId := uuid.New()
	planId, err := is.RunTerraformPlan(projectId)

	if planId == uuid.Nil {
		t.Errorf("error: %v details: %v", "planId is nil", err)
	}
}
