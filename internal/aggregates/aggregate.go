package aggregates

import "github.com/google/uuid"

type Aggregate interface {
	GetID() uuid.UUID
}
