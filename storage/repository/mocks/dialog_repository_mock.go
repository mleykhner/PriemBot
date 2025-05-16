package mocks

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type MockDialogRepository struct {
	CreateDialogFunc                 func(applicantID int64, tx *gorm.DB) (*models.Dialog, error)
	GetDialogByIDFunc                func(id uint, tx *gorm.DB) (*models.Dialog, error)
	UpdateDialogFunc                 func(dialog *models.Dialog, tx *gorm.DB) error
	GetActiveDialogByApplicantIDFunc func(applicantID int64, tx *gorm.DB) (*models.Dialog, error)
	GetActiveDialogByOperatorIDFunc  func(operatorID int64, tx *gorm.DB) (*models.Dialog, error)
	AssignOperatorFunc               func(dialogID uint, operatorID int64, tx *gorm.DB) error
	GetOpenDialogsFunc               func(tx *gorm.DB) ([]models.Dialog, error)
}

func (m *MockDialogRepository) CreateDialog(applicantID int64, tx *gorm.DB) (*models.Dialog, error) {
	return m.CreateDialogFunc(applicantID, tx)
}

func (m *MockDialogRepository) GetDialogByID(id uint, tx *gorm.DB) (*models.Dialog, error) {
	return m.GetDialogByIDFunc(id, tx)
}

func (m *MockDialogRepository) UpdateDialog(dialog *models.Dialog, tx *gorm.DB) error {
	return m.UpdateDialogFunc(dialog, tx)
}

func (m *MockDialogRepository) GetActiveDialogByApplicantID(applicantID int64, tx *gorm.DB) (*models.Dialog, error) {
	return m.GetActiveDialogByApplicantIDFunc(applicantID, tx)
}

func (m *MockDialogRepository) GetActiveDialogByOperatorID(operatorID int64, tx *gorm.DB) (*models.Dialog, error) {
	return m.GetActiveDialogByOperatorIDFunc(operatorID, tx)
}

func (m *MockDialogRepository) AssignOperator(dialogID uint, operatorID int64, tx *gorm.DB) error {
	return m.AssignOperatorFunc(dialogID, operatorID, tx)
}

func (m *MockDialogRepository) GetOpenDialogs(tx *gorm.DB) ([]models.Dialog, error) {
	return m.GetOpenDialogsFunc(tx)
}
