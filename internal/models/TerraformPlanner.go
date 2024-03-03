package models

import (
	"fmt"
	"github.com/google/uuid"
)

type ScheduledPlan bool

type Planner interface {
	Plan(planId uuid.UUID) (ScheduledPlan, error)
}

type TerraformPlanner struct {
}

func NewTerraformPlanner() (*TerraformPlanner, error) {
	return &TerraformPlanner{}, nil
}

func (e *TerraformPlanner) Plan(planId uuid.UUID) (ScheduledPlan, error) {
	fmt.Println("Hello World!")
	//todo implement scheduled plan
	return true, nil
}
