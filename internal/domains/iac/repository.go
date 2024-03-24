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

type Repository[T aggregates.Aggregate] interface {
	Get(uuid.UUID) (*T, error)
	Add(*T) error
	Update(*T) error
}
