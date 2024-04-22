package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"labraboard/internal/aggregates"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/repositories"
	iacSvc "labraboard/internal/services/iac"
	vo "labraboard/internal/valueobjects"
	"labraboard/internal/valueobjects/iacPlans"
	"os"
)

func HandlePlan(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) {
	pl := eventSubscriber.Subscribe(eb.TRIGGERED_PLAN, context.Background())
	go func(repository *repositories.UnitOfWork) {
		for msg := range pl {
			var event = events.PlanTriggered{}
			err := json.Unmarshal(msg, &event)
			if err != nil {
				fmt.Errorf("cannot handle message type %T", event)
			}
			fmt.Println("Received message:", msg)
			handlePlanTriggered(repository, event)
		}
	}(unitOfWork)
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

func handlePlanTriggered(unitOfWork *repositories.UnitOfWork, obj events.PlanTriggered) {
	iac, err := unitOfWork.IacRepository.Get(obj.ProjectId)
	if err != nil {
		panic(err)
	}
	repoUrl, repoBranch, repoPath := iac.GetRepo()
	if repoUrl == "" {
		panic("Missing repo url")
	}

	folderPath := fmt.Sprintf("/tmp/%s", obj.PlanId)
	tofuFolderPath := fmt.Sprintf("%s/%s", folderPath, repoPath)

	gitRepo, err := git.PlainClone(folderPath, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})

	branchConfig, err := gitRepo.Branch(repoBranch)
	if err != nil {
		panic(err)
	}
	defer func(folderPath string) {
		err := os.RemoveAll(folderPath)
		if err != nil {
			return
		}
	}(folderPath)

	//if err = unitOfWork.IacRepository.Update(iac); err != nil {
	//	panic(err)
	//}
	if err := createBackendFile(tofuFolderPath, "./.local-state"); err != nil {
		panic(err)
	}

	tofu, err := iacSvc.NewTofuIacService(tofuFolderPath)
	if err != nil {
		panic(err)
	}

	iacTerraformPlanJson, err := tofu.Plan(iac.GetEnvs(false), iac.GetVariables())
	if err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		if err = unitOfWork.IacRepository.Update(iac); err != nil {
			panic(err)
		}
		return
	}

	historyConfiguration := &iacPlans.HistoryProjectConfig{
		GitSha:   branchConfig.Remote,
		GitPath:  repoBranch,
		GitUrl:   repoUrl,
		Envs:     iac.GetEnvs(true),
		Variable: iac.GetVariables(),
	}

	plan, err := aggregates.NewIacPlan(obj.PlanId, aggregates.Tofu, historyConfiguration)
	if err != nil {
		panic(err) //todo handle it
	}

	plan.AddPlan(iacTerraformPlanJson.GetPlan())
	plan.AddChanges(iacTerraformPlanJson.GetChanges()...)
	iac.UpdatePlan(obj.PlanId, vo.Succeed) //optimistic change :)
	if err = unitOfWork.IacPlan.Add(plan); err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
	}
	if err = unitOfWork.IacRepository.Update(iac); err != nil {
		panic(err)
	}
}
