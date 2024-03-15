package iac

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
)

type Type string

const (
	Tofu      Type = "tofu"
	Terraform Type = "terraform"
)

type Plan struct {
	Type Type
	Id   uuid.UUID
	plan *aggregates.IacPlan
}

type LabraboardIacService interface {
	Plan(planId uuid.UUID) (*Plan, error)
}

func (plan *Plan) GetPlan() *aggregates.IacPlan {
	return plan.plan
}
