package repository

import (
	"PriemBot/storage/models"
	"errors"
	"gorm.io/gorm"
)

type DialogRepository interface {
	CreateDialog(applicantID int64, tx *gorm.DB) (*models.Dialog, error)
	GetDialogByID(id uint, tx *gorm.DB) (*models.Dialog, error)
	UpdateDialog(dialog *models.Dialog, tx *gorm.DB) error
	GetActiveDialogByApplicantID(studentID int64, tx *gorm.DB) (*models.Dialog, error)
	GetActiveDialogByOperatorID(operatorID int64, tx *gorm.DB) (*models.Dialog, error)
	AssignOperator(dialogID uint, operatorID int64, tx *gorm.DB) error
	GetOpenDialogs(tx *gorm.DB) ([]models.Dialog, error)
}

type DialogRepositoryImpl struct{}

func NewDialogRepository() DialogRepository {
	return DialogRepositoryImpl{}
}

func (d DialogRepositoryImpl) CreateDialog(applicantID int64, tx *gorm.DB) (*models.Dialog, error) {
	res := &models.Dialog{
		ApplicantID: applicantID,
		Status:      models.DialogStatusOpen,
	}
	if err := tx.Create(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (d DialogRepositoryImpl) GetDialogByID(id uint, tx *gorm.DB) (*models.Dialog, error) {
	var dialog models.Dialog
	if err := tx.Where("id = ?", id).First(&dialog).Error; err != nil {
		return nil, err
	}
	return &dialog, nil
}

func (d DialogRepositoryImpl) UpdateDialog(dialog *models.Dialog, tx *gorm.DB) error {
	return tx.Save(dialog).Error
}

func (d DialogRepositoryImpl) GetActiveDialogByApplicantID(applicantID int64, tx *gorm.DB) (*models.Dialog, error) {
	var dialog models.Dialog
	if err := tx.Where("applicant_id = ?", applicantID).Where("status = ?", models.DialogStatusActive).First(&dialog).Error; err != nil {
		return nil, err
	}
	return &dialog, nil
}

func (d DialogRepositoryImpl) GetActiveDialogByOperatorID(applicantID int64, tx *gorm.DB) (*models.Dialog, error) {
	var dialog models.Dialog
	if err := tx.Where("operator_id = ?", applicantID).Where("status = ?", models.DialogStatusActive).First(&dialog).Error; err != nil {
		return nil, err
	}
	return &dialog, nil
}

func (d DialogRepositoryImpl) AssignOperator(dialogID uint, operatorID int64, tx *gorm.DB) error {
	var dialog models.Dialog
	if err := tx.Where("id = ?", dialogID).First(&dialog).Error; err != nil {
		return err
	}

	if dialog.Status != models.DialogStatusOpen {
		return errors.New("dialog status is not open")
	}

	dialog.OperatorID = &operatorID
	dialog.Status = models.DialogStatusActive
	return tx.Save(dialog).Error
}

func (d DialogRepositoryImpl) GetOpenDialogs(tx *gorm.DB) ([]models.Dialog, error) {
	var dialogs []models.Dialog
	if err := tx.Where("status = ?", models.DialogStatusOpen).Find(&dialogs).Error; err != nil {
		return nil, err
	}
	return dialogs, nil
}
