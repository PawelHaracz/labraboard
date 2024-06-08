package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"labraboard/internal/aggregates"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"labraboard/internal/repositories"
	iacSvc "labraboard/internal/services/iac"
	vo "labraboard/internal/valueobjects"
	"labraboard/internal/valueobjects/iac"
)

// todo redesing how to treat plan aggregate to keep whole runs history
type triggeredPlanHandler struct {
	eventSubscriber  eb.EventSubscriber
	unitOfWork       *repositories.UnitOfWork
	assembler        *iacSvc.Assembler
	serviceDiscovery string
}

func newTriggeredPlanHandler(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork, discovery string) (*triggeredPlanHandler, error) {
	return &triggeredPlanHandler{
		eventSubscriber,
		unitOfWork,
		iacSvc.NewAssembler(unitOfWork),
		discovery,
	}, nil
}

func (handler *triggeredPlanHandler) Handle(ctx context.Context) {
	log := logger.GetWitContext(ctx).With().Str("event", string(events.TRIGGERED_PLAN)).Logger()
	pl := handler.eventSubscriber.Subscribe(events.TRIGGERED_PLAN, log.WithContext(ctx))
	for msg := range pl {
		var event = events.PlanTriggered{}
		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Error().Err(fmt.Errorf("cannot handle message type %T", event))
		}
		log.Trace().Msgf("Received message: %s", msg)
		err = handler.handlePlanTriggered(event, log.WithContext(ctx))
		if err != nil {
			log.Error().Err(err).Msgf("Cannot successful create plan %s", event.PlanId)
		} else {
			log.Info().Msgf("successful create the plan %s", event.PlanId)
		}
	}

}

func (handler *triggeredPlanHandler) handlePlanTriggered(obj events.PlanTriggered, ctx context.Context) error {
	log := logger.GetWitContext(ctx).With().Str("planId", obj.PlanId.String()).Str("projectId", obj.ProjectId.String()).Logger()
	var input = iacSvc.Input{
		ProjectId:    obj.ProjectId,
		PlanId:       obj.PlanId,
		Variables:    obj.Variables,
		EnvVariables: obj.EnvVariables,
		CommitName:   obj.Commit.Name,
		CommitType:   obj.Commit.Type,
		RepoPath:     obj.RepoPath,
	}
	var assembly, err = handler.assembler.Assemble(input, ctx)

	if err != nil {
		log.Error().Err(err)
		return errors.Wrap(err, "Cannot assembly of event")
	}

	folderPath := fmt.Sprintf("/tmp/%s", assembly.PlanId)
	tofuFolderPath := fmt.Sprintf("%s/%s", folderPath, assembly.RepoPath)

	git, err := iacSvc.GitClone(assembly.RepoUrl, folderPath, assembly.CommitName, assembly.CommitType)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return errors.Wrap(err, fmt.Sprintf("Cannot checkin tag %s", assembly.CommitName))
	}

	defer func(git *iacSvc.Git) {
		err = git.Clear()
		if err != nil {
			log.Error().Err(err)
			return
		}
	}(git)

	if err = createLabraboardBackendFile(tofuFolderPath, handler.serviceDiscovery, assembly.ProjectId.String()); err != nil {
		log.Error().Err(err)
		return errors.Wrap(err, "Cannot create backend")
	}

	tofu, err := iacSvc.NewTofuIacService(tofuFolderPath, log.WithContext(ctx))
	if err != nil {
		log.Error().Err(err)
		return errors.Wrap(err, "Cannot initialize tofu")
	}

	iacAggregate, err := handler.unitOfWork.IacRepository.Get(input.ProjectId, log.WithContext(ctx))
	if err != nil {
		log.Error().Err(err).Msg("missing project")
		return errors.Wrap(err, "missing project")
	}

	iacTerraformPlanJson, err := tofu.Plan(assembly.InlineEnvVariable(), assembly.InlineVariable(), log.WithContext(ctx))
	if err != nil {
		iacAggregate.UpdatePlan(obj.PlanId, vo.Failed)
		log.Warn().Err(err).Msg(err.Error())
		if err = handler.unitOfWork.IacRepository.Update(iacAggregate, log.WithContext(ctx)); err != nil {
			log.Warn().Err(err).Msg(err.Error())
			return errors.Wrap(err, "cannot update iac")
		}
		return errors.Wrap(err, "failed generate plan")
	}

	historyEnvs := make([]vo.IaCEnv, len(assembly.EnvVariables))
	if assembly.EnvVariables != nil || len(obj.EnvVariables) != 0 {
		i := 0
		for _, env := range assembly.EnvVariables {
			historyEnvs[i] = vo.IaCEnv{
				Name:      env.Name,
				Value:     env.Value,
				HasSecret: env.HasSecret,
			}
			if historyEnvs[i].HasSecret {
				historyEnvs[i].Value = vo.SECRET_VALUE_HASH
			}
		}
	}

	plan, err := handler.unitOfWork.IacPlan.Get(obj.PlanId, log.WithContext(ctx))
	if err != nil {
		historyConfiguration := &iac.HistoryProjectConfig{
			GitSha:   git.GetCommitSha(),
			GitPath:  assembly.RepoPath,
			GitUrl:   assembly.RepoUrl,
			Envs:     historyEnvs,
			Variable: assembly.Variables,
		}

		plan, err = aggregates.NewIacPlan(obj.PlanId, aggregates.Tofu, historyConfiguration)
		if err != nil {
			log.Error().Err(err)
			return errors.Wrap(err, "cannot create new plan aggregate")
		}
		if err = handler.unitOfWork.IacPlan.Add(plan, log.WithContext(ctx)); err != nil {
			log.Error().Err(err)
			return errors.Wrap(err, "cannot save plan aggregate into db")
		}
	}

	plan.AddPlan(iacTerraformPlanJson.GetPlan())
	plan.AddChanges(iacTerraformPlanJson.GetChanges()...)
	iacAggregate.UpdatePlan(obj.PlanId, vo.Succeed) //optimistic change :)
	if err = handler.unitOfWork.IacPlan.Update(plan, log.WithContext(ctx)); err != nil {
		iacAggregate.UpdatePlan(obj.PlanId, vo.Failed)
		log.Error().Err(err).Msg("Cannot update plan aggregate into db")
		return errors.Wrap(err, "Cannot update plan aggregate into db")
	}
	if err = handler.unitOfWork.IacRepository.Update(iacAggregate, log.WithContext(ctx)); err != nil {
		log.Error().Err(err)
		return errors.Wrap(err, "Cannot update iac aggregate into db")
	}
	log.Info().Msg("successful handle the event")
	return nil
}
