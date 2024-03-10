package events

import "github.com/google/uuid"

type PlanTriggered struct {
	ProjectId uuid.UUID
	PlanId    uuid.UUID
}
