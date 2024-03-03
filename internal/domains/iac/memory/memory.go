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
	iacs map[uuid.UUID]aggregates.Iac
	sync.Mutex
}

// NewRepository returns a new Repository
func NewRepository() (*Repository, error) {
	return &Repository{
		iacs: make(map[uuid.UUID]aggregates.Iac),
	}, nil
}

// Get finds a customer by ID
func (mr *Repository) Get(id uuid.UUID) (aggregates.Iac, error) {
	if iac, ok := mr.iacs[id]; ok {
		return iac, nil
	}

	return aggregates.NewIac(id)
	//return aggregates.Iac{}, iac.ErrIacNotFound
}

// Add will add a new customer to the repository
func (mr *Repository) Add(c aggregates.Iac) error {
	if mr.iacs == nil {
		// Saftey check if customers is not create, shouldn't happen if using the Factory, but you never know
		mr.Lock()
		mr.iacs = make(map[uuid.UUID]aggregates.Iac)
		mr.Unlock()
	}
	// Make sure Customer isn't already in the repository
	if _, ok := mr.iacs[c.GetID()]; ok {
		return fmt.Errorf("customer already exists: %w", iac.ErrFailedToAddIac)
	}
	mr.Lock()
	mr.iacs[c.GetID()] = c
	mr.Unlock()
	return nil
}

// Update will replace an existing customer information with the new customer information
func (mr *Repository) Update(c aggregates.Iac) error {
	// Make sure Customer is in the repository
	if _, ok := mr.iacs[c.GetID()]; !ok {
		return fmt.Errorf("customer does not exist: %w", iac.ErrUpdateIac)
	}
	mr.Lock()
	mr.iacs[c.GetID()] = c
	mr.Unlock()
	return nil
}
