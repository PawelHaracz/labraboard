package iac

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"labraboard/internal/logger"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
	vo "labraboard/internal/valueobjects"
)

var (
	EmptyIacAssembly = Output{}
)

type Output struct {
	ProjectId    uuid.UUID
	PlanId       uuid.UUID
	Variables    map[string]string
	EnvVariables []vo.IaCEnv
	CommitName   string
	CommitType   models.CommitType
	RepoPath     string
	RepoUrl      string
}

type Assembler struct {
	unitOfWork *repositories.UnitOfWork
}

type Input struct {
	ProjectId    uuid.UUID
	PlanId       uuid.UUID
	Variables    map[string]string
	EnvVariables map[string]string
	CommitName   string
	CommitType   models.CommitType
	RepoPath     string
}

func NewAssembler(unitOfWork *repositories.UnitOfWork) *Assembler {
	return &Assembler{
		unitOfWork,
	}
}

func (assembler *Assembler) Assemble(input Input, ctx context.Context) (Output, error) {
	log := logger.GetWitContext(ctx).With().Str("planId", input.PlanId.String()).Str("projectId", input.ProjectId.String()).Logger()
	iac, err := assembler.unitOfWork.IacRepository.Get(input.ProjectId, log.WithContext(ctx))

	if err != nil {
		log.Error().Err(err)
		return EmptyIacAssembly, err
	}
	repoUrl, repoBranch, repoPath := iac.GetRepo()
	eventSha, commitType, eventRepoPath := input.CommitName, input.CommitType, input.RepoPath

	plan, err := assembler.unitOfWork.IacPlan.Get(input.PlanId, log.WithContext(ctx))
	if err == nil && plan.HistoryConfig != nil {
		eventSha = plan.HistoryConfig.GitSha
		commitType = models.SHA
		eventRepoPath = plan.HistoryConfig.GitPath
	}
	if repoUrl == "" {
		err = errors.New("Missing repo url")
		log.Error().Err(err)
		return EmptyIacAssembly, err
	}
	if eventRepoPath == "" {
		eventRepoPath = repoPath
	}

	var output = Output{
		PlanId:    input.PlanId,
		ProjectId: input.ProjectId,
		RepoUrl:   repoUrl,
		RepoPath:  eventRepoPath,
	}

	if eventSha == "" {
		output.CommitType = models.BRANCH
		output.CommitName = repoBranch
	} else {
		output.CommitType = commitType
		output.CommitName = eventSha
	}
	allCurrentEnvs := iac.GetValueEnvs(false)
	var voEnvVariables []vo.IaCEnv

	for _, env := range allCurrentEnvs {
		voEnvVariables = append(voEnvVariables, env)
	}

	if plan != nil && plan.HistoryConfig != nil && len(plan.HistoryConfig.Envs) != 0 {
		for _, env := range plan.HistoryConfig.Envs {
			if env.HasSecret {
				//not reachable value so skip
				continue
			}
			var updated = false
			for i2, _ := range voEnvVariables {
				if voEnvVariables[i2].Name == env.Name {
					voEnvVariables[i2].Value = env.Value
					updated = true
				}
			}
			if !updated {
				voEnvVariables = append(voEnvVariables, env)
			}
		}
	}

	if input.EnvVariables != nil && len(input.EnvVariables) > 0 {
		for key, env := range input.EnvVariables {
			var updated = false
			for i2, env2 := range voEnvVariables {
				if env2.Name == key {
					voEnvVariables[i2].Value = env
					updated = true
					break
				}
			}
			if !updated {
				voEnvVariables = append(voEnvVariables, vo.IaCEnv{
					Name:      key,
					Value:     env,
					HasSecret: false,
				})
			}
		}
	}

	variableMap := iac.GetVariableMap()

	if plan != nil && plan.HistoryConfig != nil && len(plan.HistoryConfig.Variable) != 0 {
		for key, val := range plan.HistoryConfig.Variable {
			variableMap[key] = val
		}
	}
	if input.Variables != nil || len(input.Variables) != 0 {
		for key, val := range input.Variables {
			variableMap[key] = val
		}
	}
	output.EnvVariables = voEnvVariables
	output.Variables = variableMap

	return output, nil
}

func (output Output) InlineVariable() []string {
	var variables []string
	for key, value := range output.Variables {
		variables = append(variables, fmt.Sprintf("%s=%s", key, value))
	}
	return variables

}

func (output Output) InlineEnvVariable() map[string]string {
	var envVariables map[string]string
	for _, env := range output.EnvVariables {
		envVariables[env.Name] = env.Value
	}
	return envVariables

}
