package models

import (
	"gorm.io/gorm"
	"time"
)

type Invite struct {
	gorm.Model
	Code      string `gorm:"uniqueIndex"`
	ExpiredAt time.Time
	CreatedBy int64
	UsedBy    *int64
}
