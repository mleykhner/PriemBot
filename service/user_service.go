package service

import (
	"PriemBot/config"
	"PriemBot/storage/models"
	"PriemBot/storage/repository"
	"fmt"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(telegramID int64, username string, role models.UserRole) (*models.User, error)
	UpdateUser(user *models.User) error
	SetUserRole(telegramID int64, role models.UserRole) (*models.User, error)
	GetOperators() ([]models.User, error)
	CreateInvite(operatorID int64) (*models.Invite, error)
	ApplyInvite(telegramID int64, inviteCode string) error
	GetUserByTelegramID(telegramID int64) (*models.User, error)
	CreateInviteLink(code string) string
}

type UserServiceImpl struct {
	config     *config.BotConfig
	userRepo   repository.UserRepository
	inviteRepo repository.InviteRepository
	db         *gorm.DB
}

func NewUserService(db *gorm.DB, config *config.BotConfig) UserService {
	return &UserServiceImpl{
		config:     config,
		userRepo:   repository.NewUserRepository(),
		inviteRepo: repository.NewInviteRepository(),
		db:         db,
	}
}

func (s *UserServiceImpl) CreateUser(telegramID int64, username string, role models.UserRole) (*models.User, error) {
	return s.userRepo.CreateUser(telegramID, username, role, s.db)
}

func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	return s.userRepo.UpdateUser(user, s.db)
}

func (s *UserServiceImpl) SetUserRole(telegramID int64, role models.UserRole) (*models.User, error) {
	return s.userRepo.SetUserRoleByTelegramID(telegramID, role, s.db)
}

func (s *UserServiceImpl) GetOperators() ([]models.User, error) {
	return s.userRepo.GetOperators(s.db)
}

func (s *UserServiceImpl) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	return s.userRepo.GetUserByTelegramID(telegramID, s.db)
}

func (s *UserServiceImpl) ApplyInvite(telegramID int64, inviteCode string) error {
	// Начинаем транзакцию
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// В случае ошибки откатываем транзакцию
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Получаем пользователя
	user, err := s.userRepo.GetUserByTelegramID(telegramID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Проверяем приглашение
	if _, err := s.inviteRepo.GetInviteByCode(inviteCode, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Используем приглашение
	if err := s.inviteRepo.UseInvite(inviteCode, telegramID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Меняем роль пользователя на оператора
	user.Role = models.RoleOperator
	if err := s.userRepo.UpdateUser(user, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Если все прошло успешно, фиксируем транзакцию
	return tx.Commit().Error
}

func (s *UserServiceImpl) CreateInvite(operatorID int64) (*models.Invite, error) {
	invite, err := s.inviteRepo.CreateInvite(operatorID, s.db)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

func (s *UserServiceImpl) CreateInviteLink(code string) string {
	return fmt.Sprintf("https://t.me/%s?start=%s", s.config.Username, code)
}
