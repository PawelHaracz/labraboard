package handlers

import (
	"github.com/pkg/errors"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/helpers"
	"labraboard/internal/repositories"
)

var (
	MissingHandlerImplementedFactory = errors.New("Missing handler implemented in factory")
	HandlerAlreadyRegistered         = errors.New("Handler has already registered")
)

type EventHandlerFactory struct {
	eventSubscriber eb.EventSubscriber
	unitOfWork      *repositories.UnitOfWork
	allowedEvents   []events.EventName
}

func NewEventHandlerFactory(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) *EventHandlerFactory {
	return &EventHandlerFactory{
		eventSubscriber: eventSubscriber,
		unitOfWork:      unitOfWork,
		allowedEvents:   []events.EventName{events.LEASE_LOCK, events.TRIGGERED_PLAN},
	}
}

func (factory *EventHandlerFactory) RegisterHandler(event events.EventName) (EventHandler, error) {
	if !helpers.Contains(factory.allowedEvents, event) {
		return nil, HandlerAlreadyRegistered
	}
	switch event {
	case events.LEASE_LOCK:
		factory.allowedEvents = helpers.Remove(factory.allowedEvents, event)
		return newTerraformStateLeaseLockHandler(factory.eventSubscriber, factory.unitOfWork)
	case events.TRIGGERED_PLAN:
		factory.allowedEvents = helpers.Remove(factory.allowedEvents, event)
		return newTriggeredPlanHandler(factory.eventSubscriber, factory.unitOfWork)
	}
	return nil, MissingHandlerImplementedFactory
}

func (factory *EventHandlerFactory) RegisterAllHandlers() ([]EventHandler, error) {
	handlers := make([]EventHandler, len(factory.allowedEvents))
	for i, allowedEvent := range factory.allowedEvents {
		handler, err := factory.RegisterHandler(allowedEvent)
		if err != nil {
			return nil, err
		}
		handlers[i] = handler
	}
	return handlers, nil
}
