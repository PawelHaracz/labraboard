package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
)

type scheduledPlanHandler struct {
	eventSubscriber eb.EventSubscriber
	unitOfWork      *repositories.UnitOfWork
	publisher       eb.EventPublisher
}

func newScheduledPlanHandler(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork, publisher eb.EventPublisher) (*scheduledPlanHandler, error) {
	return &scheduledPlanHandler{
		eventSubscriber,
		unitOfWork,
		publisher,
	}, nil
}
func (handler *scheduledPlanHandler) Handle(ctx context.Context) {
	log := logger.GetWitContext(ctx).With().Str("event", events.SCHEDULED_PLAN).Logger()
	locks := handler.eventSubscriber.Subscribe(events.SCHEDULED_PLAN, ctx)
	for msg := range locks {
		var event = events.ScheduledPlan{}
		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Error().Err(fmt.Errorf("cannot handle message type %T", event))
		}
		log.Info().Msgf("Received message: %s", msg)
		go handler.handle(event, log.WithContext(ctx))
	}
}

func (handler *scheduledPlanHandler) handle(event events.ScheduledPlan, ctx context.Context) {
	log := logger.GetWitContext(ctx).
		With().
		Str("planId", event.PlanId.String()).
		Str("projectId", event.ProjectId.String()).
		Logger()

	plan, err := handler.unitOfWork.IacPlan.Get(event.PlanId, log.WithContext(ctx))
	if err != nil {
		log.Error().Err(err)
		return
	}

	if plan.HistoryConfig == nil {
		log.Error().Msg("Plan doesn't have config history")
		return
	}
	log.Info().Msg("publishing event run plan")
	handler.publisher.Publish(events.TRIGGERED_PLAN,
		events.PlanTriggered{
			ProjectId: event.ProjectId,
			PlanId:    event.PlanId,
			RepoPath:  plan.HistoryConfig.GitPath,
			Commit: events.Commit{
				Type: models.SHA,
				Name: plan.HistoryConfig.GitSha,
			},
			Variables:    plan.HistoryConfig.Variable,
			EnvVariables: plan.HistoryConfig.Envs,
		},
		log.WithContext(ctx))
	log.Info().Msg("published event run plan")
}
