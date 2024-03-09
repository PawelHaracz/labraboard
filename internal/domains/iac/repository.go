package iac

import (
	"errors"
	"labraboard/internal/aggregates"

	"github.com/google/uuid"
)

var (
	// ErrIacNotFound is returned when a IaC is not found.
	ErrIacNotFound = errors.New("the IaC was not found in the repositories")
	// ErrFailedToAddIac is returned when the IaC could not be added to the repositories.
	ErrFailedToAddIac = errors.New("failed to add the IaC to the repositories")
	// ErrUpdateIac is returned when the IaC could not be updated in the repositories.
	ErrUpdateIac = errors.New("failed to update the IaC in the repositories")
)

type Repository interface {
	Get(uuid.UUID) (aggregates.Iac, error)
	Add(aggregates.Iac) error
	Update(aggregates.Iac) error
}
