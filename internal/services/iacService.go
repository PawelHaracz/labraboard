package services

import (
	"errors"
	"github.com/google/uuid"
	"labraboard/internal/domains/iac"
	"labraboard/internal/eventbus"
)

type IacConfiguration func(os *IacService) error

type IacService struct {
	publisher  eventbus.EventPublisher
	repository iac.Repository
}

func NewIacService(configs ...IacConfiguration) (*IacService, error) {
	is := &IacService{}
	for _, cfg := range configs {
		if err := cfg(is); err != nil {
			return nil, err
		}
	}

	if is.publisher == nil {
		return nil, errors.New("planner is not set")
	}
	if is.repository == nil {
		return nil, errors.New("repositories is not set")
	}

	return is, nil
}

func WithEventBus(eb eventbus.EventPublisher) IacConfiguration {
	return func(is *IacService) error {
		is.publisher = eb
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

	svc.publisher.Publish(eventbus.TRIGGERED_PLAN, planId)
	if err != nil {
		return uuid.Nil, err
	}

	return planId, nil
}
