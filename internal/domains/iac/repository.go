package iac

import (
	"errors"
	"labraboard/internal/aggregates"

	"github.com/google/uuid"
)

var (
	// ErrIacNotFound is returned when a IaC is not found.
	ErrIacNotFound = errors.New("the IaC was not found in the repository")
	// ErrFailedToAddIac is returned when the IaC could not be added to the repository.
	ErrFailedToAddIac = errors.New("failed to add the IaC to the repository")
	// ErrUpdateIac is returned when the IaC could not be updated in the repository.
	ErrUpdateIac = errors.New("failed to update the IaC in the repository")
)

type Repository interface {
	Get(uuid.UUID) (aggregates.Iac, error)
	Add(aggregates.Iac) error
	Update(aggregates.Iac) error
}
