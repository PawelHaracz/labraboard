package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"labraboard/internal/aggregates"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/logger"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
	iacSvc "labraboard/internal/services/iac"
	vo "labraboard/internal/valueobjects"
	"labraboard/internal/valueobjects/iacPlans"
	"os"
)

// /todo redesing how to treat plan aggregate to keep whole runs history
type triggeredPlanHandler struct {
	eventSubscriber eb.EventSubscriber
	unitOfWork      *repositories.UnitOfWork
	assembler       *iacSvc.Assembler
}

func newTriggeredPlanHandler(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) (*triggeredPlanHandler, error) {
	return &triggeredPlanHandler{
		eventSubscriber,
		unitOfWork,
		iacSvc.NewAssembler(unitOfWork),
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
		log.Info().Msgf("Received message: %s", msg)
		handler.handlePlanTriggered(event, log.WithContext(ctx))
	}

}

func createBackendFile(path string, statePath string) error {
	content := `terraform {
  backend "local" {
    path = "%s"
  }
}`

	file, err := os.Create(fmt.Sprintf("%s/backend_override.tf", path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, content, statePath)
	if err != nil {
		return err
	}

	return nil
}

func (handler *triggeredPlanHandler) handlePlanTriggered(obj events.PlanTriggered, ctx context.Context) {
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
		return
	}
	folderPath := fmt.Sprintf("/tmp/%s", assembly.PlanId)
	tofuFolderPath := fmt.Sprintf("%s/%s", folderPath, assembly.RepoPath)

	gitRepo, err := git.PlainClone(folderPath, false, &git.CloneOptions{
		URL:      assembly.RepoUrl,
		Progress: os.Stdout,
	})

	defer func(folderPath string) {
		err = os.RemoveAll(folderPath)
		if err != nil {
			log.Error().Err(err)
			return
		}
	}(folderPath)
	var commitSha = ""
	switch assembly.CommitType {
	case models.TAG:
		tag, err := gitRepo.Tag(assembly.CommitName)
		if err != nil {
			log.Error().Err(err)
			return
		}
		commitSha = tag.Hash().String()
	case models.SHA:
		object, err := gitRepo.CommitObject(plumbing.NewHash(assembly.CommitName))
		if err != nil {
			log.Error().Err(err)
			return
		}
		commitSha = object.Hash.String()
	case models.BRANCH:
		branchConfig, err := gitRepo.CommitObject(plumbing.NewHash(assembly.CommitName))
		if err != nil {
			log.Error().Err(err)
			return
		}
		commitSha = branchConfig.Hash.String()
	}

	if err = createBackendFile(tofuFolderPath, "./.local-state"); err != nil {
		log.Error().Err(err)
		return
	}

	tofu, err := iacSvc.NewTofuIacService(tofuFolderPath)
	if err != nil {
		log.Error().Err(err)
		return
	}

	iacTerraformPlanJson, err := tofu.Plan(assembly.InlineEnvVariable(), assembly.InlineVariable(), log.WithContext(ctx))
	iac, err := handler.unitOfWork.IacRepository.Get(input.ProjectId, log.WithContext(ctx))
	if err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		log.Warn().Err(err).Msg(err.Error())
		if err = handler.unitOfWork.IacRepository.Update(iac, log.WithContext(ctx)); err != nil {
			log.Warn().Err(err).Msg(err.Error())
			return
		}
		return
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
		historyConfiguration := &iacPlans.HistoryProjectConfig{
			GitSha:   commitSha,
			GitPath:  assembly.RepoPath,
			GitUrl:   assembly.RepoUrl,
			Envs:     historyEnvs,
			Variable: assembly.Variables,
		}

		plan, err = aggregates.NewIacPlan(obj.PlanId, aggregates.Tofu, historyConfiguration)
		if err != nil {
			log.Error().Err(err)
			return
		}
		if err = handler.unitOfWork.IacPlan.Add(plan, log.WithContext(ctx)); err != nil {
			log.Error().Err(err)
			return
		}
	}
	plan.AddPlan(iacTerraformPlanJson.GetPlan())
	plan.AddChanges(iacTerraformPlanJson.GetChanges()...)
	iac.UpdatePlan(obj.PlanId, vo.Succeed) //optimistic change :)

	if err = handler.unitOfWork.IacPlan.Update(plan, log.WithContext(ctx)); err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		log.Error().Err(err)
		return
	}
	if err = handler.unitOfWork.IacRepository.Update(iac, log.WithContext(ctx)); err != nil {
		log.Error().Err(err)
		return
	}
}
