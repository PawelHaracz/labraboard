package main

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/repositories"
	iacSvc "labraboard/internal/services/iac"
	vo "labraboard/internal/valueobjects"
	"os"
)

func handlePlan(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) {
	pl := eventSubscriber.Subscribe(eb.TRIGGERED_PLAN, context.Background())
	go func(repository *repositories.UnitOfWork) {
		for msg := range pl {
			switch obj := msg.(type) {
			case events.PlanTriggered:
				fmt.Println("Received message:", msg)
				handlePlanTriggered(repository, obj)
			default:
				fmt.Errorf("cannot handle message type %T", obj)
			}
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
	if iac.Repo == nil {
		panic("Missing repo url")
	}

	folderPath := fmt.Sprintf("/tmp/%s", obj.PlanId)
	tofuFolderPath := fmt.Sprintf("%s/%s", folderPath, iac.Repo.Path)

	gitRepo, err := git.PlainClone(folderPath, false, &git.CloneOptions{
		URL:      iac.Repo.Url,
		Progress: os.Stdout,
	})

	if _, err := gitRepo.Branch(iac.Repo.DefaultBranch); err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(folderPath)
		if err != nil {
			return
		}
	}()

	if err = unitOfWork.IacRepository.Update(iac); err != nil {
		panic(err)
	}
	if err := createBackendFile(tofuFolderPath, "./.local-state"); err != nil {
		panic(err)
	}

	tofu, err := iacSvc.NewTofuIacService(tofuFolderPath)
	if err != nil {
		panic(err)
	}

	plan, err := tofu.Plan(obj.PlanId, iac.GetEnvs(), iac.GetVariables())
	if err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		if err = unitOfWork.IacRepository.Update(iac); err != nil {
			panic(err)
		}
	}
	iac.UpdatePlan(obj.PlanId, vo.Succeed)
	if err = unitOfWork.IacPlan.Add(plan.GetPlan()); err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		if err = unitOfWork.IacRepository.Update(iac); err != nil {
			panic(err)
		}
	}

	iac.UpdatePlan(obj.PlanId, vo.Succeed)
	if err = unitOfWork.IacRepository.Update(iac); err != nil {
		panic(err)
	}
}
