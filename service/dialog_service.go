package service

import (
	"PriemBot/storage/models"
	"PriemBot/storage/repository"

	"gorm.io/gorm"
)

type DialogsService interface {
	CreateDialog(applicantID int64) (*models.Dialog, error)

	GetDialogByID(id uint) (*models.Dialog, error)

	GetActiveDialogByStudentID(studentID int64) (*models.Dialog, error)

	GetActiveDialogByOperatorID(operatorID int64) (*models.Dialog, error)

	AssignOperator(dialogID uint, operatorID int64) error

	GetOpenDialogs() ([]models.Dialog, error)

	CreateMessage(dialogID uint, senderID int64, text string) (*models.Message, error)

	GetDialogMessages(dialogID uint) ([]models.Message, error)

	UpdateDialog(dialog *models.Dialog) error
}

type DialogsServiceImpl struct {
	dialogRepo  repository.DialogRepository
	messageRepo repository.MessageRepository
	db          *gorm.DB
}

func NewDialogsService(db *gorm.DB) DialogsService {
	return &DialogsServiceImpl{
		dialogRepo:  repository.NewDialogRepository(),
		messageRepo: repository.NewMessageRepository(),
		db:          db,
	}
}

func (s *DialogsServiceImpl) CreateDialog(applicantID int64) (*models.Dialog, error) {
	return s.dialogRepo.CreateDialog(applicantID, s.db)
}

func (s *DialogsServiceImpl) GetDialogByID(id uint) (*models.Dialog, error) {
	return s.dialogRepo.GetDialogByID(id, s.db)
}

func (s *DialogsServiceImpl) GetActiveDialogByStudentID(studentID int64) (*models.Dialog, error) {
	return s.dialogRepo.GetActiveDialogByApplicantID(studentID, s.db)
}

func (s *DialogsServiceImpl) GetActiveDialogByOperatorID(operatorID int64) (*models.Dialog, error) {
	return s.dialogRepo.GetActiveDialogByOperatorID(operatorID, s.db)
}

func (s *DialogsServiceImpl) AssignOperator(dialogID uint, operatorID int64) error {
	return s.dialogRepo.AssignOperator(dialogID, operatorID, s.db)
}

func (s *DialogsServiceImpl) GetOpenDialogs() ([]models.Dialog, error) {
	return s.dialogRepo.GetOpenDialogs(s.db)
}

func (s *DialogsServiceImpl) CreateMessage(dialogID uint, senderID int64, text string) (*models.Message, error) {
	return s.messageRepo.CreateMessage(dialogID, senderID, text, s.db)
}

func (s *DialogsServiceImpl) GetDialogMessages(dialogID uint) ([]models.Message, error) {
	return s.messageRepo.GetDialogMessages(dialogID, s.db)
}

func (s *DialogsServiceImpl) UpdateDialog(dialog *models.Dialog) error {
	return s.dialogRepo.UpdateDialog(dialog, s.db)
}
