package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IaCPlanDb struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid"`
	ChangeSummary []byte    `gorm:"null,type:jsonb"`
	Changes       []byte    `gorm:"null,type:jsonb"`
	PlanType      string    `gorm:"not null"`
	PlanJson      []byte    `gorm:"null,type:jsonb"`
	Config        []byte    `gorm:"null,type:jsonb"`
}
