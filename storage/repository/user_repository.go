package repository

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(telegramID int64, username string, role models.UserRole, tx *gorm.DB) (*models.User, error)
	UpdateUser(user *models.User, tx *gorm.DB) error
	SetUserRoleByTelegramID(id int64, role models.UserRole, tx *gorm.DB) (*models.User, error)
	GetOperators(tx *gorm.DB) ([]models.User, error)
	GetApplicants(tx *gorm.DB) ([]models.User, error)
	GetUserByTelegramID(telegramID int64, tx *gorm.DB) (*models.User, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return UserRepositoryImpl{}
}

func (u UserRepositoryImpl) CreateUser(telegramID int64, username string, role models.UserRole, tx *gorm.DB) (*models.User, error) {
	user := &models.User{
		TelegramID: telegramID,
		Name:       username,
		Role:       role,
	}
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserRepositoryImpl) UpdateUser(user *models.User, tx *gorm.DB) error {
	return tx.Save(user).Error
}

func (u UserRepositoryImpl) SetUserRoleByTelegramID(id int64, role models.UserRole, tx *gorm.DB) (*models.User, error) {
	var user models.User
	if err := tx.Where("telegram_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	user.Role = role
	if err := tx.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepositoryImpl) GetOperators(tx *gorm.DB) ([]models.User, error) {
	var operators []models.User
	if err := tx.Where("role = ?", models.RoleOperator).Find(&operators).Error; err != nil {
		return nil, err
	}
	return operators, nil
}

func (u UserRepositoryImpl) GetApplicants(tx *gorm.DB) ([]models.User, error) {
	var applicants []models.User
	if err := tx.Where("role = ?", models.RoleApplicant).Find(&applicants).Error; err != nil {
		return nil, err
	}
	return applicants, nil
}

func (u UserRepositoryImpl) GetUserByTelegramID(telegramID int64, tx *gorm.DB) (*models.User, error) {
	var user models.User
	if err := tx.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
