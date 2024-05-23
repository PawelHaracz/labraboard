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
}

func newTriggeredPlanHandler(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) (*triggeredPlanHandler, error) {
	return &triggeredPlanHandler{
		eventSubscriber,
		unitOfWork,
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
	iac, err := handler.unitOfWork.IacRepository.Get(obj.ProjectId, log.WithContext(ctx))

	if err != nil {
		log.Error().Err(err)
		return
	}

	repoUrl, repoBranch, repoPath := iac.GetRepo()
	eventSha, eventCommitType, eventRepoPath := obj.Commit.Name, obj.Commit.Type, obj.RepoPath
	if repoUrl == "" {
		log.Error().Msg("Missing repo url")
		return
	}

	if eventRepoPath == "" {
		eventRepoPath = repoPath
	}
	if eventSha == "" {
		eventSha = repoBranch

	}

	folderPath := fmt.Sprintf("/tmp/%s", obj.PlanId)
	tofuFolderPath := fmt.Sprintf("%s/%s", folderPath, eventRepoPath)

	gitRepo, err := git.PlainClone(folderPath, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})

	defer func(folderPath string) {
		err = os.RemoveAll(folderPath)
		if err != nil {
			log.Error().Err(err)
			return
		}
	}(folderPath)

	commitSha := ""

	if eventSha == "" {
		branchConfig, err := gitRepo.CommitObject(plumbing.NewHash(eventSha))
		if err != nil {
			log.Error().Err(err)
			return
		}
		commitSha = branchConfig.Hash.String()
	} else {
		switch eventCommitType {
		case models.TAG:
			tag, err := gitRepo.Tag(eventSha)
			if err != nil {
				log.Error().Err(err)
				return
			}
			commitSha = tag.Hash().String()
		case models.SHA:
			object, err := gitRepo.CommitObject(plumbing.NewHash(eventSha))
			if err != nil {
				log.Error().Err(err)
				return
			}
			commitSha = object.Hash.String()
		case models.BRANCH:
			branchConfig, err := gitRepo.CommitObject(plumbing.NewHash(eventSha))
			if err != nil {
				log.Error().Err(err)
				return
			}
			commitSha = branchConfig.Hash.String()
		}
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
	envVariables := obj.EnvVariables
	allCurrentEnvs := iac.GetEnvs(false)

	if envVariables == nil || len(envVariables) == 0 {
		envVariables = allCurrentEnvs
	}

	for key, val := range envVariables {
		if val == vo.SECRET_VALUE_HASH {
			secret, ok := allCurrentEnvs[key]
			if ok {
				envVariables[key] = secret
			}
		}
	}
	var variableMap map[string]string
	if obj.Variables != nil || len(obj.Variables) != 0 {
		variableMap = obj.Variables
	} else {
		variableMap = iac.GetVariableMap()
	}

	var variables []string
	for key, value := range variableMap {
		variables = append(variables, fmt.Sprintf("%s=%s", key, value))
	}

	iacTerraformPlanJson, err := tofu.Plan(envVariables, variables, log.WithContext(ctx))
	if err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		log.Warn().Err(err).Msg(err.Error())
		if err = handler.unitOfWork.IacRepository.Update(iac, log.WithContext(ctx)); err != nil {
			log.Warn().Err(err).Msg(err.Error())
			return
		}
		return
	}

	var envs map[string]string
	if obj.Variables != nil || len(obj.Variables) != 0 {
		envs = obj.EnvVariables
	} else {
		envs = iac.GetEnvs(true)
	}
	historyEnvs := make([]vo.IaCEnv, len(envs))
	i := 0
	for key, value := range envs {
		historyEnvs[i] = vo.IaCEnv{
			Name:      key,
			Value:     value,
			HasSecret: value == vo.SECRET_VALUE_HASH,
		}
	}

	plan, err := handler.unitOfWork.IacPlan.Get(obj.PlanId, log.WithContext(ctx))
	if err != nil {
		historyConfiguration := &iacPlans.HistoryProjectConfig{
			GitSha:   commitSha,
			GitPath:  repoBranch,
			GitUrl:   repoUrl,
			Envs:     historyEnvs,
			Variable: variableMap,
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
