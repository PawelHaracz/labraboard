package memory

import (
	"fmt"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
	"labraboard/internal/domains/iac"
	"sync"
)

// Repository fulfills the IacRepository interface
type Repository struct {
	iacs  map[uuid.UUID]*aggregates.Iac
	plans map[uuid.UUID]*aggregates.IacPlan
	sync.Mutex
	states map[uuid.UUID]*aggregates.TerraformState
}

// NewRepository returns a new Repository
func NewRepository() (*Repository, error) {
	return &Repository{
		iacs:   make(map[uuid.UUID]*aggregates.Iac),
		plans:  make(map[uuid.UUID]*aggregates.IacPlan),
		states: make(map[uuid.UUID]*aggregates.TerraformState),
	}, nil
}

// Get finds a customer by ID
func (mr *Repository) Get(id uuid.UUID) (*aggregates.Iac, error) {
	if iac, ok := mr.iacs[id]; ok {
		return iac, nil
	}

	return nil, fmt.Errorf("customer does not exist: %w", iac.ErrIacNotFound)
}

// Add will add a new customer to the repositories
func (mr *Repository) Add(c *aggregates.Iac) error {
	if mr.iacs == nil {
		// Saftey check if customers is not create, shouldn't happen if using the Factory, but you never know
		mr.Lock()
		mr.iacs = make(map[uuid.UUID]*aggregates.Iac)
		mr.Unlock()
	}
	// Make sure Customer isn't already in the repositories
	if _, ok := mr.iacs[c.GetID()]; ok {
		return fmt.Errorf("customer already exists: %w", iac.ErrFailedToAddIac)
	}
	mr.Lock()
	mr.iacs[c.GetID()] = c
	mr.Unlock()
	return nil
}

// Update will replace an existing customer information with the new customer information
func (mr *Repository) Update(c *aggregates.Iac) error {
	// Make sure Customer is in the repositories
	if _, ok := mr.iacs[c.GetID()]; !ok {
		return fmt.Errorf("customer does not exist: %w", iac.ErrUpdateIac)
	}
	mr.Lock()
	mr.iacs[c.GetID()] = c
	mr.Unlock()
	return nil
}

// GetPlan Get finds a customer by ID
func (mr *Repository) GetPlan(id uuid.UUID) (*aggregates.IacPlan, error) {
	if plan, ok := mr.plans[id]; ok {
		return plan, nil
	}

	return aggregates.NewIacPlan(id, aggregates.Terraform, nil)
}

// AddPlan Add will add a new customer to the repositories
func (mr *Repository) AddPlan(c *aggregates.IacPlan) error {
	if mr.plans == nil {
		// Saftey check if customers is not create, shouldn't happen if using the Factory, but you never know
		mr.Lock()
		mr.plans = make(map[uuid.UUID]*aggregates.IacPlan)
		mr.Unlock()
	}
	// Make sure Customer isn't already in the repositories
	if _, ok := mr.plans[c.GetID()]; ok {
		return fmt.Errorf("customer already exists: %w", iac.ErrFailedToAddIac)
	}
	mr.Lock()
	mr.plans[c.GetID()] = c
	mr.Unlock()
	return nil
}

// UpdatePlan Update will replace an existing customer information with the new customer information
func (mr *Repository) UpdatePlan(c *aggregates.IacPlan) error {
	// Make sure Customer is in the repositories
	if _, ok := mr.plans[c.GetID()]; !ok {
		return fmt.Errorf("customer does not exist: %w", iac.ErrUpdateIac)
	}
	mr.Lock()
	mr.plans[c.GetID()] = c
	mr.Unlock()
	return nil
}

func (mr *Repository) GetState(id uuid.UUID) (*aggregates.TerraformState, error) {
	if state, ok := mr.states[id]; ok {
		return state, nil
	}

	return nil, fmt.Errorf("customer does not exist: %w", iac.ErrIacNotFound)
}

func (mr *Repository) AddState(c *aggregates.TerraformState) error {
	if mr.states == nil {
		// Saftey check if customers is not create, shouldn't happen if using the Factory, but you never know
		mr.Lock()
		mr.states = make(map[uuid.UUID]*aggregates.TerraformState)
		mr.Unlock()
	}
	// Make sure Customer isn't already in the repositories
	if _, ok := mr.states[c.GetID()]; ok {
		return fmt.Errorf("Cannot update the state")
	}
	mr.Lock()
	mr.states[c.GetID()] = c
	mr.Unlock()
	return nil
}

func (mr *Repository) UpdateState(c *aggregates.TerraformState) error {
	// Make sure Customer isn't already in the repositories
	if _, ok := mr.states[c.GetID()]; !ok {
		return fmt.Errorf("Cannot update the state")
	}
	mr.Lock()
	mr.states[c.GetID()] = c
	mr.Unlock()
	return nil
}
