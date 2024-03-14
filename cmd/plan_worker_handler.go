package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	dbmemory "labraboard/internal/domains/iac/memory"
	eb "labraboard/internal/eventbus"
	"labraboard/internal/eventbus/events"
	iacSvc "labraboard/internal/services/iac"
	vo "labraboard/internal/valueobjects"
	"os"
)

func handlePlan(repository *dbmemory.Repository) {
	pl := eventBus.Subscribe(eb.TRIGGERED_PLAN)
	go func(repository *dbmemory.Repository) {
		for msg := range pl {
			switch obj := msg.(type) {
			case events.PlanTriggered:
				fmt.Println("Received message:", msg)
				handlePlanTriggered(repository, obj)
			default:
				fmt.Errorf("cannot handle message type %T", obj)
			}
		}
	}(repository)
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

func handlePlanTriggered(repository *dbmemory.Repository, obj events.PlanTriggered) {
	gitRepo, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/microsoft/terraform-azure-devops-starter.git",
		Progress: os.Stdout,
	})
	if _, err := gitRepo.Branch("master"); err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll("/tmp/foo")
		if err != nil {
			return
		}
	}()

	iac, err := repository.Get(obj.ProjectId)

	if err = repository.Update(iac); err != nil {
		panic(err)
	}
	if err := createBackendFile("/tmp/foo/101-terraform-job/terraform", "./.local-state"); err != nil {
		panic(err)
	}

	tofu, err := iacSvc.NewTofuIacService("/tmp/foo/101-terraform-job/terraform", true)
	if err != nil {
		panic(err)
	}
	envs := map[string]string{
		"ARM_TENANT_ID":       "4c83ec3e-26b4-444f-afb7-8b171cd1b420",
		"ARM_CLIENT_ID":       "99cc9476-40fd-48b6-813f-e79e0ff830fc",
		"ARM_CLIENT_SECRET":   "fixit",
		"ARM_SUBSCRIPTION_ID": "cb5863b1-784d-4813-b2c7-e87919081ecb",
	}
	plan, err := tofu.Plan(obj.PlanId, envs)
	if err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		if err = repository.Update(iac); err != nil {
			panic(err)
		}
	}
	iac.UpdatePlan(obj.PlanId, vo.Succeed)
	if err = repository.AddPlan(*plan.GetPlan()); err != nil {
		iac.UpdatePlan(obj.PlanId, vo.Failed)
		if err = repository.Update(iac); err != nil {
			panic(err)
		}
	}

	iac.UpdatePlan(obj.PlanId, vo.Succeed)
	if err = repository.Update(iac); err != nil {
		panic(err)
	}
}
