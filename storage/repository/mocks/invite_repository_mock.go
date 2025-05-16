package mocks

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type MockInviteRepository struct {
	CreateInviteFunc    func(createdBy int64, tx *gorm.DB) (*models.Invite, error)
	GetInviteByCodeFunc func(code string, tx *gorm.DB) (*models.Invite, error)
	UseInviteFunc       func(code string, usedBy int64, tx *gorm.DB) error
}

func (m *MockInviteRepository) CreateInvite(createdBy int64, tx *gorm.DB) (*models.Invite, error) {
	return m.CreateInviteFunc(createdBy, tx)
}

func (m *MockInviteRepository) GetInviteByCode(code string, tx *gorm.DB) (*models.Invite, error) {
	return m.GetInviteByCodeFunc(code, tx)
}

func (m *MockInviteRepository) UseInvite(code string, usedBy int64, tx *gorm.DB) error {
	return m.UseInviteFunc(code, usedBy, tx)
}
