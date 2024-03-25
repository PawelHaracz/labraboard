package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IaCDb struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid"`
	IacType   int       `gorm:"not null,type:integer"`
	Repo      []byte    `gorm:"null,type:jsonb"`
	Envs      []byte    `gorm:"null,type:jsonb"`
	Variables []byte    `gorm:"null,type:jsonb"`
	Plans     []byte    `gorm:"null,type:jsonb"`
}
