package mocks

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type MockUserRepository struct {
	CreateUserFunc              func(telegramID int64, username string, role models.UserRole, tx *gorm.DB) (*models.User, error)
	UpdateUserFunc              func(user *models.User, tx *gorm.DB) error
	SetUserRoleByTelegramIDFunc func(id int64, role models.UserRole, tx *gorm.DB) (*models.User, error)
	GetOperatorsFunc            func(tx *gorm.DB) ([]models.User, error)
	GetApplicantsFunc           func(tx *gorm.DB) ([]models.User, error)
	GetUserByTelegramIDFunc     func(telegramID int64, tx *gorm.DB) (*models.User, error)
}

func (m *MockUserRepository) CreateUser(telegramID int64, username string, role models.UserRole, tx *gorm.DB) (*models.User, error) {
	return m.CreateUserFunc(telegramID, username, role, tx)
}

func (m *MockUserRepository) UpdateUser(user *models.User, tx *gorm.DB) error {
	return m.UpdateUserFunc(user, tx)
}

func (m *MockUserRepository) SetUserRoleByTelegramID(id int64, role models.UserRole, tx *gorm.DB) (*models.User, error) {
	return m.SetUserRoleByTelegramIDFunc(id, role, tx)
}

func (m *MockUserRepository) GetOperators(tx *gorm.DB) ([]models.User, error) {
	return m.GetOperatorsFunc(tx)
}

func (m *MockUserRepository) GetApplicants(tx *gorm.DB) ([]models.User, error) {
	return m.GetApplicantsFunc(tx)
}

func (m *MockUserRepository) GetUserByTelegramID(telegramID int64, tx *gorm.DB) (*models.User, error) {
	return m.GetUserByTelegramIDFunc(telegramID, tx)
}
