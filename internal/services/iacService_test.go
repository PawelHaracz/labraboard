package services

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/domains/iac/memory"
	"labraboard/internal/models"
	"testing"
)

func TestNewIacService(t *testing.T) {
	tfPlanner, _ := models.NewTerraformPlanner()
	r, _ := memory.NewRepository()

	is, err := NewIacService(
		WithPlanner(tfPlanner),
		WithRepository(r))

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if is == nil {
		t.Errorf("error: %v", "IacService is nil")
	}

	if is.planner == nil {
		t.Errorf("error: %v", "IacService.planner is nil")
	}

	if is.repository == nil {
		t.Errorf("error: %v", "IacService.repository is nil")
	}

	aggregate, _ := aggregates.NewIac(uuid.New())
	err = is.repository.Add(aggregate)

	if err != nil {
		t.Errorf("error during adding item: %v", err)
	}
}

func TestIacService_RunTerraformPlan(t *testing.T) {
	tfPlanner, _ := models.NewTerraformPlanner()
	r, _ := memory.NewRepository()

	is, _ := NewIacService(
		WithPlanner(tfPlanner),
		WithRepository(r))

	projectId := uuid.New()
	planId, err := is.RunTerraformPlan(projectId)

	if planId == uuid.Nil {
		t.Errorf("error: %v details: %v", "planId is nil", err)
	}
}
