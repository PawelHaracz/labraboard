package memory

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"labraboard/internal/aggregates"
	"sync"
)

var (
	ErrNotFound = errors.New("the Item was not found in the repositories")
	// ErrFailedToAddIac is returned when the IaC could not be added to the repositories.
	ErrFailedToAddIac = errors.New("failed to add the IaC to the repositories")
	// ErrUpdateIac is returned when the IaC could not be updated in the repositories.
	ErrUpdateIac = errors.New("failed to update the IaC in the repositories")
)

type GenericRepository[T aggregates.Aggregate] struct {
	collection map[uuid.UUID]T
	sync.Mutex
}

func NewGenericRepository[T aggregates.Aggregate]() *GenericRepository[T] {
	return &GenericRepository[T]{
		collection: make(map[uuid.UUID]T),
	}
}

func (r *GenericRepository[T]) Get(id uuid.UUID, ctx context.Context) (T, error) {
	if iac, ok := r.collection[id]; ok {
		return iac, nil
	}
	return getZero[T](), fmt.Errorf("Not found: %w", ErrNotFound)
}

func (r *GenericRepository[T]) Add(t T, ctx context.Context) error {
	if r.collection == nil {
		// Saftey check if customers is not create, shouldn't happen if using the Factory, but you never know
		r.Lock()
		r.collection = make(map[uuid.UUID]T)
		r.Unlock()
	}
	// Make sure Customer isn't already in the repositories
	if _, ok := r.collection[t.GetID()]; ok {
		return fmt.Errorf("Not found: %w", ErrNotFound)
	}
	r.Lock()
	r.collection[t.GetID()] = t
	r.Unlock()
	return nil
}

func (r *GenericRepository[T]) Update(t T, ctx context.Context) error {
	// Make sure Customer is in the repositories
	if _, ok := r.collection[t.GetID()]; !ok {
		return fmt.Errorf("Not found: %w", ErrNotFound)
	}
	r.Lock()
	r.collection[t.GetID()] = t
	r.Unlock()
	return nil
}

func (r *GenericRepository[T]) GetAll(ctx context.Context) []T {
	var collection []T
	for _, item := range r.collection {
		collection = append(collection, item)
	}
	return collection
}

func getZero[T any]() T {
	var result T
	return result
}
