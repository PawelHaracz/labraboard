package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TerraformStateDb struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid"`
	State     []byte    `gorm:"null,type:bytea"`
	CreatedOn time.Time `gorm:"default:CURRENT_TIMESTAMP()`
	ModifyOn  time.Time `gorm:"default:CURRENT_TIMESTAMP()`
	Lock      []byte    `gorm:"null,type:bytea"`
}
