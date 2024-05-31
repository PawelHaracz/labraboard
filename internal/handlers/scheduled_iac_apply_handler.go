package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"labraboard/internal/repositories"
	"labraboard/internal/services/iac"
)

type scheduledIaCApplyHandler struct {
	eventSubscriber eb.EventSubscriber
	unitOfWork      *repositories.UnitOfWork
}

func newScheduledIaCApplyHandler(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) (*scheduledIaCApplyHandler, error) {
	return &scheduledIaCApplyHandler{
		eventSubscriber,
		unitOfWork,
	}, nil
}

func (handler *scheduledIaCApplyHandler) Handle(ctx context.Context) {
	log := logger.GetWitContext(ctx).With().Str("event", string(events.IAC_APPLY_SCHEDULED)).Logger()
	locks := handler.eventSubscriber.Subscribe(events.SCHEDULED_PLAN, log.WithContext(ctx))
	for msg := range locks {
		var event = events.IacApplyScheduled{}
		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Error().Err(fmt.Errorf("cannot handle message type %T", event))
		}
		log.Info().Msgf("Received message: %s", msg)
		go handler.handle(event, log.WithContext(ctx))
	}
}

func (handler *scheduledIaCApplyHandler) handle(event events.IacApplyScheduled, ctx context.Context) {
	log := logger.GetWitContext(ctx).
		With().
		Str("changeId", event.ChangeId.String()).
		Str("iacType", string(event.IacType)).
		Str("planId", event.PlanId.String()).
		Str("projectId", event.ProjectId.String()).
		Str("owner", event.Owner).
		Logger()
	//_, err := handler.unitOfWork.IacPlan.Get(event.PlanId, log.WithContext(ctx))

	assembler := iac.NewAssembler(handler.unitOfWork)

	var input = iac.Input{
		ProjectId:    event.ProjectId,
		PlanId:       event.PlanId,
		Variables:    nil,
		EnvVariables: nil,
		CommitName:   "",
		CommitType:   "",
		RepoPath:     "",
	}

	_, err := assembler.Assemble(input, log.WithContext(ctx))
	if err != nil {
		log.Error().Err(err)
		return
	}
	//save tfplan to tfplan
	//run apply

}
