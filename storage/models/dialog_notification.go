package models

import "gorm.io/gorm"

type DialogNotification struct {
	gorm.Model
	TelegramUserID    int64
	DialogID          uint `gorm:"index"`
	TelegramMessageID int
}
