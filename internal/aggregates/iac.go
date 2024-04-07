package aggregates

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"labraboard/internal/repositories/postgres/models"
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
	Repo      *vo.IaCRepo
	variables []*vo.IaCVariable
}

func NewIac(id uuid.UUID, iacType vo.IaCType, plans []*vo.Plans, envs []*vo.IaCEnv, repo *vo.IaCRepo, variables []*vo.IaCVariable) (*Iac, error) {
	aggregate := &Iac{}
	aggregate.id = id
	aggregate.plans = plans
	aggregate.envs = envs
	aggregate.IacType = iacType
	aggregate.Repo = repo
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
	if receiver.Repo != nil {
		return errors.New("repo already exists")
	}
	repo, err := vo.NewIaCRepo(url, defaultBranch, path)
	receiver.Repo = repo
	return err
}

// GetID returns the Iac root entity ID
func (receiver *Iac) GetID() uuid.UUID {
	return receiver.id
}

func (receiver *Iac) AddPlan(id uuid.UUID) {
	utc := time.Now().UTC()
	receiver.plans = append(receiver.plans, &vo.Plans{
		Status:    vo.Scheduled,
		ModifyOn:  utc,
		CreatedOn: utc,
		Id:        id,
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

func (receiver *Iac) GetEnvs() map[string]string {
	var envs = map[string]string{}
	for _, env := range receiver.envs {
		envs[env.Name] = env.Value
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

func (receiver *Iac) SetVariable(name string, value string) error {
	receiver.variables = append(receiver.variables, &vo.IaCVariable{
		Name:  name,
		Value: value,
	})
	return nil
}

func (receiver *Iac) Map() (*models.IaCDb, error) {
	iacRepo, err := json.Marshal(receiver.Repo)

	if err != nil {
		return nil, errors.Wrap(err, "can't create repo on receiver")
	}

	envs, err := json.Marshal(receiver.envs)
	if err != nil {
		return nil, errors.Wrap(err, "can't create envs on receiver")
	}

	variables, err := json.Marshal(receiver.variables)
	if err != nil {
		return nil, errors.Wrap(err, "can't create variables on receiver")
	}

	plans, err := json.Marshal(receiver.plans)
	if err != nil {
		return nil, errors.Wrap(err, "can't create plans on receiver")
	}

	return &models.IaCDb{
		ID:        receiver.id,
		IacType:   int(receiver.IacType),
		Repo:      iacRepo,
		Envs:      envs,
		Variables: variables,
		Plans:     plans,
	}, nil
}
