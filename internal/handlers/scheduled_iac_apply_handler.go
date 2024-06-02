package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"labraboard/internal/repositories"
	"labraboard/internal/services/iac"
	"os"
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

func (handler *scheduledIaCApplyHandler) handle(event events.IacApplyScheduled, ctx context.Context) error {
	const tfPlanPath = "plan.tfplan"
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

	output, err := assembler.Assemble(input, log.WithContext(ctx))
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	if len(output.PlanRaw) == 0 {
		err = errors.New("Missing plan")
		log.Error().Err(err)
		return nil
	}

	folderPath := fmt.Sprintf("/tmp/%s/apply", output.PlanId)
	tofuFolderPath := fmt.Sprintf("%s/%s", folderPath, output.RepoPath)

	planPath := fmt.Sprintf("%s/%s", tofuFolderPath, tfPlanPath)
	git, err := iac.GitClone(output.RepoUrl, folderPath, output.CommitName, output.CommitType)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return errors.Wrap(err, fmt.Sprintf("Cannot checkin tag %s", output.CommitName))
	}

	defer func(git *iac.Git) {
		err = git.Clear()
		if err != nil {
			log.Error().Err(err)
			return
		}
	}(git)

	if err = createBackendFile(tofuFolderPath, "./.local-state"); err != nil {
		log.Error().Err(err)
		return nil
	}

	if err = handler.savePlanAsTfPlan(planPath, output.PlanRaw); err != nil {
		log.Error().Err(err)
		return nil
	}

	tofu, err := iac.NewTofuIacService(tofuFolderPath)
	if err != nil {
		log.Error().Err(err)
		return nil
	}
	_, err = tofu.Apply(output.PlanId, output.InlineEnvVariable(), output.InlineVariable(), planPath, ctx)

	return nil
}

func (handler *scheduledIaCApplyHandler) savePlanAsTfPlan(path string, planRaw []byte) error {
	fo, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err = fo.Close(); err != nil {
			err = errors.Wrap(err, "problem with close file")
		}
	}()

	w := bufio.NewWriter(fo)
	if _, err = w.Write(planRaw); err != nil {
		err = errors.Wrap(err, "problem with write file")
	}

	return err
}
