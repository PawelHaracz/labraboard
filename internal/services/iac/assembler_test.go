package iac

import (
	"labraboard/internal/aggregates"
	"labraboard/internal/logger"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
	dbmemory "labraboard/internal/repositories/memory"
	"labraboard/internal/valueobjects"
	"labraboard/internal/valueobjects/iac"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestAssembler_Assemble(t *testing.T) {
	logger.Init(7, false) //disabled
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.Iac](),
		),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.IacPlan](),
		),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.TerraformState](),
		),
		repositories.WithIacPlanRepositoryDbRepositoryMemory(
			dbmemory.NewGenericRepository[*aggregates.IacDeployment](),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	projectId, planIds := Arrange(uow)
	var assembler = NewAssembler(uow)
	var ctx = context.TODO()
	t.Run("Plan doesn't exist, should take from project values", func(t *testing.T) {
		var input = Input{
			ProjectId:    projectId,
			PlanId:       planIds[0],
			Variables:    nil,
			EnvVariables: nil,
			CommitName:   "",
			CommitType:   "",
			RepoPath:     "",
		}

		var expectedEnvVariables = []valueobjects.IaCEnv{
			{
				Name:      "TEST",
				Value:     "test1",
				HasSecret: false,
			},
			{
				Name:      "TEST1",
				Value:     "t",
				HasSecret: true,
			},
		}
		var expectedVariables = map[string]string{
			"VARTEST": "BLABLE",
		}

		if output, err := assembler.Assemble(input, ctx); err != nil {
			t.Fatal(err)
		} else {

			assert.Equal(t, planIds[0], output.PlanId)
			assert.Equal(t, projectId, output.ProjectId)
			assert.Equal(t, "main", output.CommitName)
			assert.Equal(t, models.BRANCH, output.CommitType)
			assert.Equal(t, "https://github.com/pawelharacz/labraboard", output.RepoUrl)
			assert.Equal(t, "", output.RepoPath)
			assert.ElementsMatch(t, expectedEnvVariables, output.EnvVariables)
			assert.EqualValues(t, expectedVariables, output.Variables)
		}
	})

	t.Run("Plan exists, but it doesn't have envs and variables should combine values from project and plan", func(t *testing.T) {
		var input = Input{
			ProjectId:    projectId,
			PlanId:       planIds[1],
			Variables:    nil,
			EnvVariables: nil,
			CommitName:   "88864e896674402e4b54e0b8aa53b77aa18fb8dd",
			CommitType:   models.SHA,
			RepoPath:     "",
		}

		var expectedEnvVariables = []valueobjects.IaCEnv{
			{
				Name:      "TEST",
				Value:     "test1",
				HasSecret: false,
			},
			{
				Name:      "TEST1",
				Value:     "t",
				HasSecret: true,
			},
		}
		var expectedVariables = map[string]string{
			"VARTEST": "BLABLE",
		}

		if output, err := assembler.Assemble(input, ctx); err != nil {
			t.Fatal(err)
		} else {

			assert.Equal(t, planIds[1], output.PlanId)
			assert.Equal(t, projectId, output.ProjectId)
			assert.Equal(t, "88864e896674402e4b54e0b8aa53b77aa18fb8dd", output.CommitName)
			assert.Equal(t, models.SHA, output.CommitType)
			assert.Equal(t, "https://github.com/pawelharacz/labraboard", output.RepoUrl)
			assert.Equal(t, "", output.RepoPath)
			assert.ElementsMatch(t, expectedEnvVariables, output.EnvVariables)
			assert.EqualValues(t, expectedVariables, output.Variables)
		}
	})

	t.Run("Plan exists, input contains provided env and variables but plan doesn't have envs and variables should combine values from input, project and plan", func(t *testing.T) {
		var input = Input{
			ProjectId: projectId,
			PlanId:    planIds[1],
			Variables: map[string]string{
				"RETEST": "LOL",
			},
			EnvVariables: map[string]string{
				"TEST":  "test2",
				"4TEST": "test4",
			},
			CommitName: "88864e896674402e4b54e0b8aa53b77aa18fb8dd",
			CommitType: models.SHA,
			RepoPath:   "",
		}

		var expectedEnvVariables = []valueobjects.IaCEnv{
			{
				Name:      "TEST",
				Value:     "test2",
				HasSecret: false,
			},
			{
				Name:      "TEST1",
				Value:     "t",
				HasSecret: true,
			},
			{
				Name:      "4TEST",
				Value:     "test4",
				HasSecret: false,
			},
		}
		var expectedVariables = map[string]string{
			"VARTEST": "BLABLE",
			"RETEST":  "LOL",
		}

		if output, err := assembler.Assemble(input, ctx); err != nil {
			t.Fatal(err)
		} else {

			assert.Equal(t, planIds[1], output.PlanId)
			assert.Equal(t, projectId, output.ProjectId)
			assert.Equal(t, "88864e896674402e4b54e0b8aa53b77aa18fb8dd", output.CommitName)
			assert.Equal(t, models.SHA, output.CommitType)
			assert.Equal(t, "https://github.com/pawelharacz/labraboard", output.RepoUrl)
			assert.Equal(t, "", output.RepoPath)
			assert.ElementsMatch(t, expectedEnvVariables, output.EnvVariables)
			assert.EqualValues(t, expectedVariables, output.Variables)
		}
	})

	t.Run("Plan exists, input contains provided env and variables, plan has envs and variables should combine values from input, project and plan", func(t *testing.T) {
		var input = Input{
			ProjectId: projectId,
			PlanId:    planIds[2],
			Variables: map[string]string{
				"RETEST": "LOL",
			},
			EnvVariables: map[string]string{
				"4TEST": "test4",
			},
			CommitName: "88864e896674402e4b54e0b8aa53b77aa18fb8dd",
			CommitType: models.SHA,
			RepoPath:   "",
		}

		var expectedEnvVariables = []valueobjects.IaCEnv{
			{
				Name:      "TEST",
				Value:     "test2",
				HasSecret: false,
			},
			{
				Name:      "TEST1",
				Value:     "t",
				HasSecret: true,
			},
			{
				Name:      "4TEST",
				Value:     "test4",
				HasSecret: false,
			},
		}
		var expectedVariables = map[string]string{
			"VARTEST": "FOOBAR",
			"RETEST":  "LOL",
		}

		if output, err := assembler.Assemble(input, ctx); err != nil {
			t.Fatal(err)
		} else {

			assert.Equal(t, planIds[2], output.PlanId)
			assert.Equal(t, projectId, output.ProjectId)
			assert.Equal(t, "88864e896674402e4b54e0b8aa53b77aa18fb8dd", output.CommitName)
			assert.Equal(t, models.SHA, output.CommitType)
			assert.Equal(t, "https://github.com/pawelharacz/labraboard", output.RepoUrl)
			assert.Equal(t, "", output.RepoPath)
			assert.ElementsMatch(t, expectedEnvVariables, output.EnvVariables)
			assert.EqualValues(t, expectedVariables, output.Variables)
		}
	})
}

func TestAssembler_Inline(t *testing.T) {
	output := Output{
		ProjectId: uuid.New(),
		PlanId:    uuid.New(),
		Variables: map[string]string{
			"TEST1": "FOO",
			"BAR":   "RAB",
		},
		EnvVariables: []valueobjects.IaCEnv{
			{
				Name:      "HERO",
				Value:     "FLASH",
				HasSecret: false,
			},
			{
				Name:      "BODY",
				Value:     "head",
				HasSecret: true,
			},
		},
		CommitName: "",
		CommitType: "",
		RepoPath:   "",
		RepoUrl:    "",
		PlanRaw:    nil,
	}

	t.Run("inline_variables", func(t *testing.T) {
		variables := output.InlineVariable()

		assert.Equal(t, 2, len(variables))
		assert.Equal(t, "TEST1=FOO", variables[0])
		assert.Equal(t, "BAR=RAB", variables[1])
	})

	t.Run("inlinve_env", func(t *testing.T) {
		var expected = map[string]string{
			"HERO": "FLASH",
			"BODY": "head",
		}

		envs := output.InlineEnvVariable()

		assert.Equal(t, 2, len(envs))
		assert.EqualValues(t, expected, envs)
	})
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
			Value:     "t",
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
		{
			Id:        planIds[2],
			Status:    valueobjects.Pending,
			CreatedOn: time.Now(),
			ModifyOn:  time.Now(),
		},
	}
	aggregateIac, _ := aggregates.NewIac(uuid.New(), valueobjects.Terraform, plans, envs, repo, variables)
	uow.IacRepository.Add(aggregateIac, context.TODO())

	var historyConfig = &iac.HistoryProjectConfig{
		GitSha:   "88864e896674402e4b54e0b8aa53b77aa18fb8dd",
		GitUrl:   "https://github.com/PawelHaracz/labraboard/",
		GitPath:  "",
		Envs:     nil,
		Variable: nil,
	}

	plan, _ := aggregates.NewIacPlan(planIds[1], aggregates.Terraform, historyConfig)
	uow.IacPlan.Add(plan, context.TODO())

	var historyConfig1 = &iac.HistoryProjectConfig{
		GitSha:  "88864e896674402e4b54e0b8aa53b77aa18fb8dd",
		GitUrl:  "https://github.com/PawelHaracz/labraboard/",
		GitPath: "",
		Envs: []valueobjects.IaCEnv{
			{
				Name:      "TEST",
				Value:     "test2",
				HasSecret: false,
			},
			{
				Name:      "TEST1",
				Value:     valueobjects.SECRET_VALUE_HASH,
				HasSecret: true,
			},
		},
		Variable: map[string]string{
			"VARTEST": "FOOBAR",
		},
	}
	plan1, _ := aggregates.NewIacPlan(planIds[2], aggregates.Terraform, historyConfig1)
	uow.IacPlan.Add(plan1, context.TODO())
	return aggregateIac.GetID(), planIds
}
