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
	eventSubscriber  eb.EventSubscriber
	unitOfWork       *repositories.UnitOfWork
	allowedEvents    []events.EventName
	eventPublisher   eb.EventPublisher
	serviceDiscovery string
}

func NewEventHandlerFactory(eventSubscriber eb.EventSubscriber, eventPublisher eb.EventPublisher, unitOfWork *repositories.UnitOfWork, serviceDiscovery string) *EventHandlerFactory {
	return &EventHandlerFactory{
		eventSubscriber:  eventSubscriber,
		unitOfWork:       unitOfWork,
		eventPublisher:   eventPublisher,
		serviceDiscovery: serviceDiscovery,
		allowedEvents:    []events.EventName{events.LEASE_LOCK, events.TRIGGERED_PLAN, events.SCHEDULED_PLAN, events.IAC_APPLY_SCHEDULED},
	}
}

func (factory *EventHandlerFactory) RegisterHandler(event events.EventName) (EventHandler, error) {
	if !helpers.Contains(factory.allowedEvents, event) {
		return nil, errors.Wrap(HandlerAlreadyRegistered, string(event))
	}
	switch event {
	case events.LEASE_LOCK:
		factory.allowedEvents = helpers.Remove(factory.allowedEvents, event)
		return newTerraformStateLeaseLockHandler(factory.eventSubscriber, factory.unitOfWork)
	case events.TRIGGERED_PLAN:
		factory.allowedEvents = helpers.Remove(factory.allowedEvents, event)
		return newTriggeredPlanHandler(factory.eventSubscriber, factory.unitOfWork, factory.serviceDiscovery)
	case events.SCHEDULED_PLAN:
		factory.allowedEvents = helpers.Remove(factory.allowedEvents, event)
		return newScheduledPlanHandler(factory.eventSubscriber, factory.unitOfWork, factory.eventPublisher)
	case events.IAC_APPLY_SCHEDULED:
		factory.allowedEvents = helpers.Remove(factory.allowedEvents, event)
		return newScheduledIaCApplyHandler(factory.eventSubscriber, factory.unitOfWork, factory.serviceDiscovery)
	}
	return nil, MissingHandlerImplementedFactory
}

func (factory *EventHandlerFactory) RegisterAllHandlers() ([]EventHandler, error) {
	handlers := make([]EventHandler, len(factory.allowedEvents))
	allowedEvents := make([]events.EventName, len(factory.allowedEvents))
	copy(allowedEvents, factory.allowedEvents)
	for i, allowedEvent := range allowedEvents {
		handler, err := factory.RegisterHandler(allowedEvent)
		if err != nil {
			return nil, err
		}
		handlers[i] = handler
	}
	return handlers, nil
}
