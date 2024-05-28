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
	allCurrentEnvs := iac.GetEnvs(false)
	var voEnvVariables = make([]vo.IaCEnv, len(allCurrentEnvs))
	if plan == nil || plan.HistoryConfig == nil || len(plan.HistoryConfig.Envs) == 0 {
		envVariables := input.EnvVariables

		if envVariables == nil || len(envVariables) == 0 {
			envVariables = allCurrentEnvs
		}
		var i = 0
		for key, val := range envVariables {
			if val == vo.SECRET_VALUE_HASH {
				secret, ok := allCurrentEnvs[key]
				if ok {
					voEnvVariables[i] = vo.IaCEnv{
						Name:      key,
						Value:     secret,
						HasSecret: true,
					}
				}
			} else {
				voEnvVariables[i] = vo.IaCEnv{
					Name:      key,
					Value:     val,
					HasSecret: false,
				}
			}
			i = i + 1
		}
	} else {
		voEnvVariables = plan.HistoryConfig.Envs
		for index, env := range voEnvVariables {
			if voEnvVariables[index].HasSecret {
				voEnvVariables[index].Value = allCurrentEnvs[env.Name]
			}
		}
	}

	var variableMap map[string]string
	if plan == nil || plan.HistoryConfig == nil || len(plan.HistoryConfig.Variable) == 0 {
		if input.Variables != nil || len(input.Variables) != 0 {
			variableMap = input.Variables
		} else {
			variableMap = iac.GetVariableMap()
		}
	} else {
		variableMap = plan.HistoryConfig.Variable
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
