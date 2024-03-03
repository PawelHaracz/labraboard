package services

import (
	"errors"
	"github.com/google/uuid"
	"labraboard/internal/domains/iac"
	"labraboard/internal/models"
)

type IacConfiguration func(os *IacService) error

type IacService struct {
	planner    models.Planner
	repository iac.Repository
}

func NewIacService(configs ...IacConfiguration) (*IacService, error) {
	is := &IacService{}
	for _, cfg := range configs {
		if err := cfg(is); err != nil {
			return nil, err
		}
	}

	if is.planner == nil {
		return nil, errors.New("planner is not set")
	}
	if is.repository == nil {
		return nil, errors.New("repository is not set")
	}

	return is, nil
}

func WithPlanner(p models.Planner) IacConfiguration {
	return func(is *IacService) error {
		is.planner = p
		return nil
	}
}

func WithRepository(r iac.Repository) IacConfiguration {
	return func(is *IacService) error {
		is.repository = r
		return nil
	}
}

func (svc *IacService) RunTerraformPlan(projectId uuid.UUID) (uuid.UUID, error) {
	planId := uuid.New()

	iac, err := svc.repository.Get(projectId)
	if err != nil {
		return uuid.Nil, err
	}

	iac.AddPlan(planId)
	err = svc.repository.Add(iac)
	if err != nil {
		return uuid.Nil, err
	}

	plan, err := svc.planner.Plan(planId)
	if plan == false {
		return uuid.Nil, err
	}

	return planId, nil
}
