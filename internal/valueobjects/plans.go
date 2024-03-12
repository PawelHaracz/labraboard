package valueobjects

import (
	"github.com/google/uuid"
	"time"
	_ "time"
)

type IaCType int

type PlanStatus string

const (
	Terraform IaCType = iota
)

const (
	Scheduled PlanStatus = "scheduled"
	Pending   PlanStatus = "pending"
	Succeed   PlanStatus = "Succeed"
	Failed    PlanStatus = "failed"
)

type Plans struct {
	Id        uuid.UUID
	Status    PlanStatus
	CreatedOn time.Time
	ModifyOn  time.Time
}
