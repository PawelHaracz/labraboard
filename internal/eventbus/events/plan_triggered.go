package events

import (
	"encoding/json"
	"github.com/google/uuid"
)

type PlanTriggered struct {
	ProjectId uuid.UUID
	PlanId    uuid.UUID
}

func (i PlanTriggered) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
