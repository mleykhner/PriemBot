package repository

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type DialogNotificationRepository interface {
	DeleteAllByDialogID(dialogID uint, tx *gorm.DB) error
	GetAllByDialogID(dialogID uint, tx *gorm.DB) ([]models.DialogNotification, error)
	DeleteByID(id uint, tx *gorm.DB) error
	CreateNotification(dialogID uint, telegramMessageID int, tx *gorm.DB) error
}

type DialogNotificationRepositoryImpl struct{}

func NewDialogNotificationRepository() DialogNotificationRepository {
	return DialogNotificationRepositoryImpl{}
}

func (d DialogNotificationRepositoryImpl) DeleteAllByDialogID(dialogID uint, tx *gorm.DB) error {
	return tx.Where("dialog_id = ?", dialogID).Delete(&models.DialogNotification{}).Error
}

func (d DialogNotificationRepositoryImpl) GetAllByDialogID(dialogID uint, tx *gorm.DB) ([]models.DialogNotification, error) {
	var dialogNotifications []models.DialogNotification
	if err := tx.Where("dialog_id = ?", dialogID).Find(&dialogNotifications).Error; err != nil {
		return nil, err
	}
	return dialogNotifications, nil
}

func (d DialogNotificationRepositoryImpl) DeleteByID(id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).Delete(&models.DialogNotification{}).Error
}

func (d DialogNotificationRepositoryImpl) CreateNotification(dialogID uint, telegramMessageID int, tx *gorm.DB) error {
	newNotification := models.DialogNotification{
		DialogID:          dialogID,
		TelegramMessageID: telegramMessageID,
	}
	return tx.Create(&newNotification).Error
}
