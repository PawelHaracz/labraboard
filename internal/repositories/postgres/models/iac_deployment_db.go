package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IaCDeploymentDb struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid"`
	PlanId         uuid.UUID `gorm:"type:uuid"`
	ProjectId      uuid.UUID `gorm:"type:uuid"`
	Started        time.Time `gorm:"type:time"`
	Deployed       time.Time `gorm:"type:time"`
	DeploymentType string    `gorm:"not null"`
	ChangeSummary  []byte    `gorm:"null,type:jsonb"`
	Changes        []byte    `gorm:"null,type:jsonb"`
	Outputs        []byte    `gorm:"null,type:jsonb"`
}
