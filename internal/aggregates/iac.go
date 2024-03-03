package aggregates

import (
	"errors"
	"github.com/google/uuid"
	vo "labraboard/internal/valueobjects"
	"time"
)

var (
	//ErrPlanNotFound is returned when a plan is not found
	ErrPlanNotFound = errors.New("plan not found")
)

type Iac struct {
	id    uuid.UUID
	plans []vo.Plans
}

func NewIac(id uuid.UUID) (Iac, error) {
	aggregate := &Iac{}
	aggregate.id = id
	aggregate.plans = []vo.Plans{}

	return *aggregate, nil
}

// GetID returns the Iac root entity ID
func (c *Iac) GetID() uuid.UUID {
	return c.id
}

func (c *Iac) AddPlan(id uuid.UUID) {
	utc := time.Now().UTC()
	c.plans = append(c.plans, vo.Plans{
		Status:    vo.Scheduled,
		ModifyOn:  utc,
		CreatedOn: utc,
		Id:        id,
		PlanType:  vo.Terraform,
	})
}

func (c *Iac) GetPlan(id uuid.UUID) (vo.Plans, error) {
	for _, plan := range c.plans {
		if plan.Id == id {
			return plan, nil
		}
	}
	return vo.Plans{}, ErrPlanNotFound
}
