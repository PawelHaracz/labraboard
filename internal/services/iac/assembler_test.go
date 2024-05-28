package iac

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"labraboard/internal/aggregates"
	"labraboard/internal/logger"
	"labraboard/internal/repositories"
	dbmemory "labraboard/internal/repositories/memory"
	"labraboard/internal/valueobjects"
	"labraboard/internal/valueobjects/iacPlans"
	"testing"
	"time"
)

func TestAssembler_Assemble(t *testing.T) {
	t.SkipNow()           //todo
	logger.Init(7, false) //disabled
	uow, _ := repositories.NewUnitOfWork(
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.IacPlan](),
		),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.Iac](),
		),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.TerraformState](),
		),
	)
	//projectId, planIds := Arrange(uow)
	_ = NewAssembler(uow)

}

func Arrange(uow *repositories.UnitOfWork) (uuid.UUID, []uuid.UUID) {
	var envs = []*valueobjects.IaCEnv{
		{
			Name:      "TEST",
			Value:     "test1",
			HasSecret: false,
		},
		{
			Name:      "TEST1",
			Value:     "test1",
			HasSecret: true,
		},
	}

	var repo = &valueobjects.IaCRepo{
		Url:           "https://github.com/pawelharacz/labraboard",
		DefaultBranch: "main",
		Path:          "",
	}

	var variables = []*valueobjects.IaCVariable{
		{
			Name:  "VARTEST",
			Value: "BLABLE",
		},
	}

	var planIds = []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	var plans = []*valueobjects.Plans{
		{
			Id:        planIds[0],
			Status:    valueobjects.Pending,
			CreatedOn: time.Now(),
			ModifyOn:  time.Now(),
		},
		{
			Id:        planIds[1],
			Status:    valueobjects.Pending,
			CreatedOn: time.Now(),
			ModifyOn:  time.Now(),
		},
	}
	iac, _ := aggregates.NewIac(uuid.New(), valueobjects.Terraform, plans, envs, repo, variables)
	uow.IacRepository.Add(iac, context.TODO())
	//todo add envs and variables
	var historyConfig = &iacPlans.HistoryProjectConfig{
		GitSha:   "88864e896674402e4b54e0b8aa53b77aa18fb8dd",
		GitUrl:   "https://github.com/PawelHaracz/labraboard/",
		GitPath:  "",
		Envs:     nil,
		Variable: nil,
	}

	plan, _ := aggregates.NewIacPlan(planIds[1], aggregates.Terraform, historyConfig)
	uow.IacPlan.Add(plan, context.TODO())
	return iac.GetID(), planIds
}
