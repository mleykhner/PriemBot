package models

import "time"

type User struct {
	TelegramID int64 `gorm:"primaryKey"`
	Name       string
	Role       UserRole `gorm:"type:varchar(20)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
