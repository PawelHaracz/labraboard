package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TerraformStateDb struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid"`
	State     []byte    `gorm:"null,type:jsonb"`
	CreatedOn time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ModifyOn  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Lock      []byte    `gorm:"null,type:jsonb"`
}
