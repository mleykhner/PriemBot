package models

import "gorm.io/gorm"

type Dialog struct {
	gorm.Model
	ApplicantID int64
	OperatorID  *int64       // Может быть nil, пока не назначен
	Status      DialogStatus `gorm:"type:varchar(20)"`
}
