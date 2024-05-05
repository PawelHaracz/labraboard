package handlers

import (
	eb "labraboard/internal/eventbus"
	"labraboard/internal/logger"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
)

//todo implement generic one handler with object instead multiple
import (
	"context"
	"encoding/json"
	"fmt"
	"labraboard/internal/eventbus/events"
)

type terraformStateLeaseLockHandler struct {
	eventSubscriber eb.EventSubscriber
	unitOfWork      *repositories.UnitOfWork
}

func newTerraformStateLeaseLockHandler(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) (*terraformStateLeaseLockHandler, error) {
	return &terraformStateLeaseLockHandler{
		eventSubscriber,
		unitOfWork,
	}, nil
}

func (handler *terraformStateLeaseLockHandler) Handle(ctx context.Context) {
	log := logger.GetWitContext(ctx).With().Str("event", string(events.LEASE_LOCK)).Logger()
	locks := handler.eventSubscriber.Subscribe(events.LEASE_LOCK, ctx)
	for msg := range locks {
		var event = events.LeasedLock{}
		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Error().Err(fmt.Errorf("cannot handle message type %T", event))
		}
		log.Info().Msgf("Received message: %s", msg)
		go handler.handle(event, log.WithContext(ctx))
	}
}

func (handler *terraformStateLeaseLockHandler) handle(event events.LeasedLock, ctx context.Context) {
	log := logger.GetWitContext(ctx).With().Str("eventType", string(event.Type)).Str("eventId", event.Id.String()).Logger()
	if event.Type != models.Terraform {
		log.Warn().Msg("wrong event type")
		return
	}
	item, err := handler.unitOfWork.TerraformStateDbRepository.Get(event.Id, log.WithContext(ctx))
	if err != nil {
		log.Error().Err(err)
	}

	info, err := item.GetLockInfo()
	if err != nil {
		log.Error().Err(err)
		return
	}

	if info == nil {
		log.Warn().Msg("empty info lock")
		return
	}

	if err = item.SetLockInfo(nil); err != nil {
		log.Error().Err(err)
		return
	}
}
