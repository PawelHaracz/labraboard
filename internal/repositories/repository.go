package repositories

import (
	"errors"
	"golang.org/x/net/context"
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
	Get(uuid.UUID, context.Context) (T, error)
	Add(T, context.Context) error
	Update(T, context.Context) error
	GetAll(context.Context) []T
}
