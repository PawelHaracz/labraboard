package aggregates

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	vo "labraboard/internal/valueobjects"
	_ "slices"
	"time"
)

var (
	//ErrPlanNotFound is returned when a plan is not found
	ErrPlanNotFound     = errors.New("plan not found")
	ErrEnvAlreadyExists = errors.New("env already exists")
)

type Iac struct {
	id        uuid.UUID
	plans     []*vo.Plans
	IacType   vo.IaCType
	envs      []*vo.IaCEnv
	repo      *vo.IaCRepo
	variables []*vo.IaCVariable
}

func NewIac(id uuid.UUID, iacType vo.IaCType, plans []*vo.Plans, envs []*vo.IaCEnv, repo *vo.IaCRepo, variables []*vo.IaCVariable) (*Iac, error) {
	aggregate := &Iac{}
	aggregate.id = id
	aggregate.plans = plans
	aggregate.envs = envs
	aggregate.IacType = iacType
	aggregate.repo = repo
	aggregate.variables = variables
	return aggregate, nil
}

func (receiver *Iac) AddEnv(name string, value string, hasSecret bool) error {
	for _, a := range receiver.envs {
		if a.Name == name {
			return ErrEnvAlreadyExists
		}
	}
	receiver.envs = append(receiver.envs, &vo.IaCEnv{
		Name:      name,
		Value:     value,
		HasSecret: hasSecret,
	})
	return nil
}

func (receiver *Iac) AddRepo(url string, defaultBranch string, path string) error {
	if receiver.repo != nil {
		return errors.New("repo already exists")
	}
	repo, err := vo.NewIaCRepo(url, defaultBranch, path)
	receiver.repo = repo
	return err
}

func (receiver *Iac) GetRepo() (Url string, DefaultBranch string, Path string) {
	return receiver.repo.Url, receiver.repo.DefaultBranch, receiver.repo.Path
}

// GetID returns the Iac root entity ID
func (receiver *Iac) GetID() uuid.UUID {
	return receiver.id
}

func (receiver *Iac) AddPlan(id uuid.UUID, sha string, path string, variables map[string]string) {
	utc := time.Now().UTC()
	receiver.plans = append(receiver.plans, &vo.Plans{
		Status:    vo.Scheduled,
		ModifyOn:  utc,
		CreatedOn: utc,
		Id:        id,
		CommitSha: sha,
		RepoPath:  path,
		Variables: variables,
	})
}

func (receiver *Iac) GetPlan(id uuid.UUID) (*vo.Plans, error) {
	for _, plan := range receiver.plans {
		if plan.Id == id {
			return plan, nil
		}
	}
	return nil, ErrPlanNotFound
}

func (receiver *Iac) GetPlans() []*vo.Plans {
	return receiver.plans
}

func (receiver *Iac) UpdatePlan(id uuid.UUID, status vo.PlanStatus) {
	if plan, err := receiver.GetPlan(id); err == nil {
		utc := time.Now().UTC()
		plan.Status = status
		plan.ModifyOn = utc
	}
}

func (receiver *Iac) GetEnvs(hideSecret bool) map[string]string {
	var envs = map[string]string{}
	for _, env := range receiver.envs {
		if hideSecret && env.HasSecret {
			envs[env.Name] = "***"
		} else {
			envs[env.Name] = env.Value
		}
	}
	return envs
}

func (receiver *Iac) GetVariables() []string {
	var variables []string
	for _, variable := range receiver.variables {
		variables = append(variables, fmt.Sprintf("%s=%s", variable.Name, variable.Value))
	}
	return variables
}

func (receiver *Iac) GetVariableMap() map[string]string {
	var variables = map[string]string{}
	for _, variable := range receiver.variables {
		variables[variable.Name] = variable.Value
	}
	return variables
}

func (receiver *Iac) SetVariable(name string, value string) error {
	receiver.variables = append(receiver.variables, &vo.IaCVariable{
		Name:  name,
		Value: value,
	})
	return nil
}

func (receiver *Iac) Composite() ([]*vo.IaCEnv, []*vo.IaCVariable, *vo.IaCRepo) {
	return receiver.envs, receiver.variables, receiver.repo
}

func (receiver *Iac) RemoveEnv(name string) error {
	for index, env := range receiver.envs {
		if env.Name == name {
			receiver.envs = append(receiver.envs[:index], receiver.envs[index+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("env %s not found", name))
}

func (receiver *Iac) RemoveVariable(name string) error {
	for index, variable := range receiver.variables {
		if variable.Name == name {
			receiver.variables = append(receiver.variables[:index], receiver.variables[index+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("variable %s not found", name))
}
