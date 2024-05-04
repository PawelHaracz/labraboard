package iac

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"labraboard/internal/aggregates"
	"labraboard/internal/models"
)

type Plan struct {
	Type models.Type
	Id   uuid.UUID
	plan *aggregates.IacPlan
}

type LabraboardIacService interface {
	Plan(planId uuid.UUID, ctx context.Context) (*Plan, error) //todo remove it
}

func (plan *Plan) GetPlan() *aggregates.IacPlan {
	return plan.plan
}
