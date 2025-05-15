package repository

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(dialogID uint, senderID int64, text string, tx *gorm.DB) (*models.Message, error)
	GetDialogMessages(dialogID uint, tx *gorm.DB) ([]models.Message, error)
}

type MessageRepositoryImpl struct{}

func NewMessageRepository() MessageRepository {
	return MessageRepositoryImpl{}
}

func (m MessageRepositoryImpl) CreateMessage(dialogID uint, senderID int64, text string, tx *gorm.DB) (*models.Message, error) {
	message := &models.Message{
		DialogID: dialogID,
		SenderID: senderID,
		Text:     text,
	}
	if err := tx.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (m MessageRepositoryImpl) GetDialogMessages(dialogID uint, tx *gorm.DB) ([]models.Message, error) {
	var messages []models.Message
	if err := tx.Where("dialog_id = ?", dialogID).Order("created_at asc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
