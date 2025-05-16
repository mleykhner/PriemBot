package tests

import (
	"PriemBot/config"
	"PriemBot/service"
	"PriemBot/storage/models"
	"PriemBot/storage/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitService_CheckAndCreateOperator(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*UserServiceMock, *mocks.MockInviteRepository)
		wantErr   bool
	}{
		{
			name: "operators exist",
			mockSetup: func(userService *UserServiceMock, inviteRepo *mocks.MockInviteRepository) {
				userService.GetOperatorsFunc = func() ([]models.User, error) {
					return []models.User{
						{TelegramID: 123, Role: models.RoleOperator},
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "no operators, create invite",
			mockSetup: func(userService *UserServiceMock, inviteRepo *mocks.MockInviteRepository) {
				userService.GetOperatorsFunc = func() ([]models.User, error) {
					return []models.User{}, nil
				}
				inviteRepo.CreateInviteFunc = func(createdBy int64, tx *gorm.DB) (*models.Invite, error) {
					return &models.Invite{
						Code:      "test-code",
						CreatedBy: createdBy,
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "error getting operators",
			mockSetup: func(userService *UserServiceMock, inviteRepo *mocks.MockInviteRepository) {
				userService.GetOperatorsFunc = func() ([]models.User, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "error creating invite",
			mockSetup: func(userService *UserServiceMock, inviteRepo *mocks.MockInviteRepository) {
				userService.GetOperatorsFunc = func() ([]models.User, error) {
					return []models.User{}, nil
				}
				inviteRepo.CreateInviteFunc = func(createdBy int64, tx *gorm.DB) (*models.Invite, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := &UserServiceMock{}
			mockInviteRepo := &mocks.MockInviteRepository{}
			tt.mockSetup(mockUserService, mockInviteRepo)

			initService := service.NewInitService(mockUserService, mockInviteRepo, nil)
			err := initService.CheckAndCreateOperator(&config.BotConfig{})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
