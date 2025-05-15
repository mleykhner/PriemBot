package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	DialogID uint
	SenderID int64
	Text     string
}
