package tests

import (
	"PriemBot/storage/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name       string
		telegramID int64
		username   string
		role       models.UserRole
		mockSetup  func(*UserServiceMock)
		wantErr    bool
	}{
		{
			name:       "successful creation",
			telegramID: 123,
			username:   "test_user",
			role:       models.RoleApplicant,
			mockSetup: func(m *UserServiceMock) {
				m.CreateUserFunc = func(telegramID int64, username string, role models.UserRole) (*models.User, error) {
					return &models.User{
						TelegramID: telegramID,
						Name:       username,
						Role:       role,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name:       "repository error",
			telegramID: 123,
			username:   "test_user",
			role:       models.RoleApplicant,
			mockSetup: func(m *UserServiceMock) {
				m.CreateUserFunc = func(telegramID int64, username string, role models.UserRole) (*models.User, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &UserServiceMock{}
			tt.mockSetup(mockService)

			user, err := mockService.CreateUser(tt.telegramID, tt.username, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.telegramID, user.TelegramID)
				assert.Equal(t, tt.username, user.Name)
				assert.Equal(t, tt.role, user.Role)
			}
		})
	}
}

func TestUserService_SetUserRole(t *testing.T) {
	tests := []struct {
		name       string
		telegramID int64
		role       models.UserRole
		mockSetup  func(*UserServiceMock)
		wantErr    bool
	}{
		{
			name:       "successful role update",
			telegramID: 123,
			role:       models.RoleOperator,
			mockSetup: func(m *UserServiceMock) {
				m.SetUserRoleFunc = func(telegramID int64, role models.UserRole) (*models.User, error) {
					return &models.User{
						TelegramID: telegramID,
						Role:       role,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name:       "user not found",
			telegramID: 123,
			role:       models.RoleOperator,
			mockSetup: func(m *UserServiceMock) {
				m.SetUserRoleFunc = func(telegramID int64, role models.UserRole) (*models.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &UserServiceMock{}
			tt.mockSetup(mockService)

			user, err := mockService.SetUserRole(tt.telegramID, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.telegramID, user.TelegramID)
				assert.Equal(t, tt.role, user.Role)
			}
		})
	}
}

func TestUserService_GetOperators(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*UserServiceMock)
		wantErr   bool
		wantCount int
	}{
		{
			name: "successful get operators",
			mockSetup: func(m *UserServiceMock) {
				m.GetOperatorsFunc = func() ([]models.User, error) {
					return []models.User{
						{TelegramID: 1, Role: models.RoleOperator},
						{TelegramID: 2, Role: models.RoleOperator},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "empty operators list",
			mockSetup: func(m *UserServiceMock) {
				m.GetOperatorsFunc = func() ([]models.User, error) {
					return []models.User{}, nil
				}
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name: "repository error",
			mockSetup: func(m *UserServiceMock) {
				m.GetOperatorsFunc = func() ([]models.User, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &UserServiceMock{}
			tt.mockSetup(mockService)

			operators, err := mockService.GetOperators()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, operators)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, operators)
				assert.Equal(t, tt.wantCount, len(operators))
			}
		})
	}
}
