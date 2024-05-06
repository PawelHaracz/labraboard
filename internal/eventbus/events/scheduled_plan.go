package events

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

const SCHEDULED_PLAN = "scheduled_plan"

type ScheduledPlan struct {
	ProjectId uuid.UUID
	PlanId    uuid.UUID
	When      time.Time
}

func (i ScheduledPlan) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
