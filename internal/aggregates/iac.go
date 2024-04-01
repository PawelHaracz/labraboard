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

func (c *Iac) AddEnv(name string, value string, hasSecret bool) error {
	for _, a := range c.envs {
		if a.Name == name {
			return ErrEnvAlreadyExists
		}
	}
	c.envs = append(c.envs, &vo.IaCEnv{
		Name:      name,
		Value:     value,
		HasSecret: hasSecret,
	})
	return nil
}

func (c *Iac) AddRepo(url string, defaultBranch string, path string) error {
	if c.Repo != nil {
		return errors.New("repo already exists")
	}
	repo, err := vo.NewIaCRepo(url, defaultBranch, path)
	c.Repo = repo
	return err
}

// GetID returns the Iac root entity ID
func (c *Iac) GetID() uuid.UUID {
	return c.id
}

func (c *Iac) AddPlan(id uuid.UUID) {
	utc := time.Now().UTC()
	c.plans = append(c.plans, &vo.Plans{
		Status:    vo.Scheduled,
		ModifyOn:  utc,
		CreatedOn: utc,
		Id:        id,
	})
}

func (c *Iac) GetPlan(id uuid.UUID) (*vo.Plans, error) {
	for _, plan := range c.plans {
		if plan.Id == id {
			return plan, nil
		}
	}
	return nil, ErrPlanNotFound
}

func (c *Iac) UpdatePlan(id uuid.UUID, status vo.PlanStatus) {
	if plan, err := c.GetPlan(id); err == nil {
		utc := time.Now().UTC()
		plan.Status = status
		plan.ModifyOn = utc
	}
}

func (iac *Iac) GetEnvs() map[string]string {
	var envs = map[string]string{}
	for _, env := range iac.envs {
		envs[env.Name] = env.Value
	}
	return envs
}

func (iac *Iac) GetVariables() []string {
	var variables []string
	for _, variable := range iac.variables {
		variables = append(variables, fmt.Sprintf("%s=%s", variable.Name, variable.Value))
	}
	return variables
}

func (iac *Iac) SetVariable(name string, value string) error {
	iac.variables = append(iac.variables, &vo.IaCVariable{
		Name:  name,
		Value: value,
	})
	return nil
}

func (iac *Iac) Map() (*models.IaCDb, error) {
	iacRepo, err := json.Marshal(iac.Repo)

	if err != nil {
		return nil, errors.Wrap(err, "can't create repo on iac")
	}

	envs, err := json.Marshal(iac.envs)
	if err != nil {
		return nil, errors.Wrap(err, "can't create envs on iac")
	}

	variables, err := json.Marshal(iac.variables)
	if err != nil {
		return nil, errors.Wrap(err, "can't create variables on iac")
	}

	plans, err := json.Marshal(iac.plans)
	if err != nil {
		return nil, errors.Wrap(err, "can't create plans on iac")
	}

	return &models.IaCDb{
		ID:        iac.id,
		IacType:   int(iac.IacType),
		Repo:      iacRepo,
		Envs:      envs,
		Variables: variables,
		Plans:     plans,
	}, nil
}
