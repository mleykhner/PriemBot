package mocks

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type MockMessageRepository struct {
	CreateMessageFunc     func(dialogID uint, senderID int64, text string, tx *gorm.DB) (*models.Message, error)
	GetDialogMessagesFunc func(dialogID uint, tx *gorm.DB) ([]models.Message, error)
}

func (m *MockMessageRepository) CreateMessage(dialogID uint, senderID int64, text string, tx *gorm.DB) (*models.Message, error) {
	return m.CreateMessageFunc(dialogID, senderID, text, tx)
}

func (m *MockMessageRepository) GetDialogMessages(dialogID uint, tx *gorm.DB) ([]models.Message, error) {
	return m.GetDialogMessagesFunc(dialogID, tx)
}
