package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TerraformStateDb struct {
	gorm.Model
	ProjectId uuid.UUID
	State     []byte
	CreatedOn time.Time
	ModifyOn  time.Time
	Lock      []byte
}
