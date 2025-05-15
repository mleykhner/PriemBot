package repository

import (
	"PriemBot/storage/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InviteRepository interface {
	CreateInvite(createdBy int64, tx *gorm.DB) (*models.Invite, error)
	GetInviteByCode(code string, tx *gorm.DB) (*models.Invite, error)
	UseInvite(code string, usedBy int64, tx *gorm.DB) error
}

type InviteRepositoryImpl struct{}

func NewInviteRepository() InviteRepository {
	return InviteRepositoryImpl{}
}

func (i InviteRepositoryImpl) CreateInvite(createdBy int64, tx *gorm.DB) (*models.Invite, error) {
	invite := &models.Invite{
		Code:      uuid.New().String(),
		CreatedBy: createdBy,
		ExpiredAt: time.Now().Add(24 * time.Hour), // Приглашение действует 24 часа
	}
	if err := tx.Create(invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (i InviteRepositoryImpl) GetInviteByCode(code string, tx *gorm.DB) (*models.Invite, error) {
	var invite models.Invite
	if err := tx.Where("code = ? AND expired_at > ?", code, time.Now()).First(&invite).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}

func (i InviteRepositoryImpl) UseInvite(code string, usedBy int64, tx *gorm.DB) error {
	var invite models.Invite
	if err := tx.Where("code = ? AND expired_at > ? AND used_by IS NULL", code, time.Now()).First(&invite).Error; err != nil {
		return err
	}
	invite.UsedBy = &usedBy
	return tx.Save(&invite).Error
}
