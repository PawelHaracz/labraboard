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
	EMPTY_IAC_ASSEMBLY = Output{}
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
		return EMPTY_IAC_ASSEMBLY, err
	}

	plan, err := assembler.unitOfWork.IacPlan.Get(input.PlanId, log.WithContext(ctx))
	if err == nil {
		log.Info().Msg(plan.GetPlanJson()) //todo handle plan changes, commits etc.
	}
	repoUrl, repoBranch, repoPath := iac.GetRepo()
	eventSha, commitType, eventRepoPath := input.CommitName, input.CommitType, input.RepoPath
	if repoUrl == "" {
		err = errors.New("Missing repo url")
		log.Error().Err(err)
		return EMPTY_IAC_ASSEMBLY, err
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

	envVariables := input.EnvVariables
	allCurrentEnvs := iac.GetEnvs(false)
	var voEnvVariables = make([]vo.IaCEnv, len(allCurrentEnvs))
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

	var variableMap map[string]string
	if input.Variables != nil || len(input.Variables) != 0 {
		variableMap = input.Variables
	} else {
		variableMap = iac.GetVariableMap()
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
