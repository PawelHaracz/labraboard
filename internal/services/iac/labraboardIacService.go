package iac

import (
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/models"
)

type Plan struct {
	Type models.Type
	Id   uuid.UUID
	plan *aggregates.IacPlan
}

type LabraboardIacService interface {
	Plan(planId uuid.UUID) (*Plan, error)
}

func (plan *Plan) GetPlan() *aggregates.IacPlan {
	return plan.plan
}
