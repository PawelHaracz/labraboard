package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"labraboard/internal/models"
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

func (handler *scheduledIaCApplyHandler) handle(event events.IacApplyScheduled, ctx context.Context) {
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
		return
	}

	if len(output.PlanRaw) == 0 {
		err = errors.New("Missing plan")
		log.Error().Err(err)
		return
	}

	folderPath := fmt.Sprintf("/tmp/%s/apply", output.PlanId)
	tofuFolderPath := fmt.Sprintf("%s/%s", folderPath, output.RepoPath)

	planPath := fmt.Sprintf("%s/%s", tofuFolderPath, tfPlanPath)
	gitRepo, err := git.PlainClone(folderPath, false, &git.CloneOptions{
		URL:      output.RepoUrl,
		Progress: os.Stdout,
	})

	defer func(folderPath string) {
		err = os.RemoveAll(folderPath)
		if err != nil {
			log.Error().Err(err)
			return
		}
	}(folderPath)

	switch output.CommitType {
	case models.TAG:
		_, err = gitRepo.Tag(output.CommitName)
		if err != nil {
			log.Error().Err(err)
			return
		}
	case models.SHA:
		_, err = gitRepo.CommitObject(plumbing.NewHash(output.CommitName))
		if err != nil {
			log.Error().Err(err)
			return
		}
	case models.BRANCH:
		_, err = gitRepo.CommitObject(plumbing.NewHash(output.CommitName))
		if err != nil {
			log.Error().Err(err)
			return
		}
	}

	if err = createBackendFile(tofuFolderPath, "./.local-state"); err != nil {
		log.Error().Err(err)
		return
	}

	if err = handler.savePlanAsTfPlan(planPath, output.PlanRaw); err != nil {
		log.Error().Err(err)
		return
	}

	tofu, err := iac.NewTofuIacService(tofuFolderPath)
	if err != nil {
		log.Error().Err(err)
		return
	}
	_, err = tofu.Apply(output.PlanId, output.InlineEnvVariable(), output.InlineVariable(), planPath, ctx)

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
