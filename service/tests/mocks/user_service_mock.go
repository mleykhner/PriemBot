package service_mocks

import (
	"PriemBot/storage/models"
)

type UserServiceMock struct {
	CreateUserFunc          func(telegramID int64, username string, role models.UserRole) (*models.User, error)
	UpdateUserFunc          func(user *models.User) error
	SetUserRoleFunc         func(telegramID int64, role models.UserRole) (*models.User, error)
	GetOperatorsFunc        func() ([]models.User, error)
	CreateInviteFunc        func(operatorID int64) (*models.Invite, error)
	ApplyInviteFunc         func(telegramID int64, inviteCode string) error
	GetUserByTelegramIDFunc func(telegramID int64) (*models.User, error)
	CreateInviteLinkFunc    func(code string) string
}

func (m *UserServiceMock) CreateUser(telegramID int64, username string, role models.UserRole) (*models.User, error) {
	return m.CreateUserFunc(telegramID, username, role)
}

func (m *UserServiceMock) UpdateUser(user *models.User) error {
	return m.UpdateUserFunc(user)
}

func (m *UserServiceMock) SetUserRole(telegramID int64, role models.UserRole) (*models.User, error) {
	return m.SetUserRoleFunc(telegramID, role)
}

func (m *UserServiceMock) GetOperators() ([]models.User, error) {
	return m.GetOperatorsFunc()
}

func (m *UserServiceMock) CreateInvite(operatorID int64) (*models.Invite, error) {
	return m.CreateInviteFunc(operatorID)
}

func (m *UserServiceMock) ApplyInvite(telegramID int64, inviteCode string) error {
	return m.ApplyInviteFunc(telegramID, inviteCode)
}

func (m *UserServiceMock) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	return m.GetUserByTelegramIDFunc(telegramID)
}

func (m *UserServiceMock) CreateInviteLink(code string) string {
	return m.CreateInviteLinkFunc(code)
}
