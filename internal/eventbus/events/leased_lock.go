package events

import (
	"encoding/json"
	"github.com/google/uuid"
	"labraboard/internal/models"
	"time"
)

const LEASE_LOCK EventName = "lease_lock"

type LeasedLock struct {
	Id        uuid.UUID
	Type      models.Type
	LeaseTime time.Time
}

func (i LeasedLock) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
