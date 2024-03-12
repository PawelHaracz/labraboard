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

	tofu, err := iacSvc.NewTofuIacService("/tmp/foo/101-terraform-job/terraform", true)
	if err != nil {
		panic(err)
	}
	plan, err := tofu.Plan(obj.PlanId)
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
