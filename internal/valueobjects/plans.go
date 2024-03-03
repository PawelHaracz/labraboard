package valueobjects

import (
	"github.com/google/uuid"
	"time"
	_ "time"
)

type PlanType int

type PlanStatus string

const (
	Terraform PlanType = iota
)

const (
	Scheduled PlanStatus = "scheduled"
	Pending   PlanStatus = "pending"
	Succeed   PlanStatus = "Succeed"
	Failed    PlanStatus = "failed"
)

type Plans struct {
	Id        uuid.UUID
	PlanType  PlanType
	Status    PlanStatus
	CreatedOn time.Time
	ModifyOn  time.Time
}
