package tests

import (
	"PriemBot/service"
	"PriemBot/storage/models"
	"PriemBot/storage/repository/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDialogService_CreateDialog(t *testing.T) {
	tests := []struct {
		name        string
		applicantID int64
		mockSetup   func(*mocks.MockDialogRepository)
		wantErr     bool
	}{
		{
			name:        "successful creation",
			applicantID: 123,
			mockSetup: func(m *mocks.MockDialogRepository) {
				m.CreateDialogFunc = func(applicantID int64, tx *gorm.DB) (*models.Dialog, error) {
					return &models.Dialog{
						ApplicantID: applicantID,
						Status:      models.DialogStatusOpen,
						Model: gorm.Model{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name:        "repository error",
			applicantID: 123,
			mockSetup: func(m *mocks.MockDialogRepository) {
				m.CreateDialogFunc = func(applicantID int64, tx *gorm.DB) (*models.Dialog, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDialogRepo := &mocks.MockDialogRepository{}
			mockMessageRepo := &mocks.MockMessageRepository{}
			tt.mockSetup(mockDialogRepo)

			dialogService := service.NewDialogsService(mockDialogRepo, mockMessageRepo, nil)
			dialog, err := dialogService.CreateDialog(tt.applicantID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dialog)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dialog)
				assert.Equal(t, tt.applicantID, dialog.ApplicantID)
				assert.Equal(t, models.DialogStatusOpen, dialog.Status)
			}
		})
	}
}

func TestDialogService_AssignOperator(t *testing.T) {
	tests := []struct {
		name       string
		dialogID   uint
		operatorID int64
		mockSetup  func(*mocks.MockDialogRepository)
		wantErr    bool
	}{
		{
			name:       "successful assignment",
			dialogID:   1,
			operatorID: 123,
			mockSetup: func(m *mocks.MockDialogRepository) {
				m.AssignOperatorFunc = func(dialogID uint, operatorID int64, tx *gorm.DB) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:       "dialog not found",
			dialogID:   1,
			operatorID: 123,
			mockSetup: func(m *mocks.MockDialogRepository) {
				m.AssignOperatorFunc = func(dialogID uint, operatorID int64, tx *gorm.DB) error {
					return errors.New("dialog not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDialogRepo := &mocks.MockDialogRepository{}
			mockMessageRepo := &mocks.MockMessageRepository{}
			tt.mockSetup(mockDialogRepo)

			dialogService := service.NewDialogsService(mockDialogRepo, mockMessageRepo, nil)
			err := dialogService.AssignOperator(tt.dialogID, tt.operatorID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDialogService_CreateMessage(t *testing.T) {
	tests := []struct {
		name      string
		dialogID  uint
		senderID  int64
		text      string
		mockSetup func(*mocks.MockMessageRepository)
		wantErr   bool
	}{
		{
			name:     "successful message creation",
			dialogID: 1,
			senderID: 123,
			text:     "test message",
			mockSetup: func(m *mocks.MockMessageRepository) {
				m.CreateMessageFunc = func(dialogID uint, senderID int64, text string, tx *gorm.DB) (*models.Message, error) {
					return &models.Message{
						DialogID: dialogID,
						SenderID: senderID,
						Text:     text,
						Model: gorm.Model{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name:     "repository error",
			dialogID: 1,
			senderID: 123,
			text:     "test message",
			mockSetup: func(m *mocks.MockMessageRepository) {
				m.CreateMessageFunc = func(dialogID uint, senderID int64, text string, tx *gorm.DB) (*models.Message, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDialogRepo := &mocks.MockDialogRepository{}
			mockMessageRepo := &mocks.MockMessageRepository{}
			tt.mockSetup(mockMessageRepo)

			dialogService := service.NewDialogsService(mockDialogRepo, mockMessageRepo, nil)
			message, err := dialogService.CreateMessage(tt.dialogID, tt.senderID, tt.text)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, message)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, message)
				assert.Equal(t, tt.dialogID, message.DialogID)
				assert.Equal(t, tt.senderID, message.SenderID)
				assert.Equal(t, tt.text, message.Text)
			}
		})
	}
}

func TestDialogService_GetDialogMessages(t *testing.T) {
	tests := []struct {
		name      string
		dialogID  uint
		mockSetup func(*mocks.MockMessageRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name:     "successful get messages",
			dialogID: 1,
			mockSetup: func(m *mocks.MockMessageRepository) {
				m.GetDialogMessagesFunc = func(dialogID uint, tx *gorm.DB) ([]models.Message, error) {
					return []models.Message{
						{DialogID: dialogID, Text: "message 1"},
						{DialogID: dialogID, Text: "message 2"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:     "empty messages list",
			dialogID: 1,
			mockSetup: func(m *mocks.MockMessageRepository) {
				m.GetDialogMessagesFunc = func(dialogID uint, tx *gorm.DB) ([]models.Message, error) {
					return []models.Message{}, nil
				}
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:     "repository error",
			dialogID: 1,
			mockSetup: func(m *mocks.MockMessageRepository) {
				m.GetDialogMessagesFunc = func(dialogID uint, tx *gorm.DB) ([]models.Message, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDialogRepo := &mocks.MockDialogRepository{}
			mockMessageRepo := &mocks.MockMessageRepository{}
			tt.mockSetup(mockMessageRepo)

			dialogService := service.NewDialogsService(mockDialogRepo, mockMessageRepo, nil)
			messages, err := dialogService.GetDialogMessages(tt.dialogID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, messages)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, messages)
				assert.Equal(t, tt.wantCount, len(messages))
			}
		})
	}
}
