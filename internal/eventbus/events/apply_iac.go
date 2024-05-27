package events

import (
	"encoding/json"
	"github.com/google/uuid"
	"labraboard/internal/aggregates"
)

const IAC_APPLY_SCHEDULED EventName = "iac_apply_scheduled"

type IacApplyScheduled struct {
	ChangeId  uuid.UUID
	ProjectId uuid.UUID
	PlanId    uuid.UUID
	IacType   aggregates.IaCPlanType
	Owner     string
}

func (i IacApplyScheduled) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
