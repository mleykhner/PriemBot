package service

import (
	"PriemBot/config"
	"PriemBot/storage/repository"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type InitService interface {
	CheckAndCreateOperator(botConfig *config.BotConfig) error
}

type InitServiceImpl struct {
	userService UserService
	inviteRepo  repository.InviteRepository
	db          *gorm.DB
}

func NewInitService(userService UserService, inviteRepo repository.InviteRepository, db *gorm.DB) InitService {
	return &InitServiceImpl{
		userService: userService,
		inviteRepo:  inviteRepo,
		db:          db,
	}
}

func (s *InitServiceImpl) CheckAndCreateOperator(botConfig *config.BotConfig) error {
	// Проверяем наличие операторов
	operators, err := s.userService.GetOperators()
	if err != nil {
		return fmt.Errorf("failed to get operators: %w", err)
	}

	// Если операторы есть, ничего не делаем
	if len(operators) > 0 {
		return nil
	}

	// Создаем приглашение
	invite, err := s.inviteRepo.CreateInvite(0, s.db) // 0 как ID создателя, так как это системное приглашение
	if err != nil {
		return fmt.Errorf("failed to create invite: %w", err)
	}

	// Формируем ссылку
	inviteLink := fmt.Sprintf("https://t.me/%s?start=%s", botConfig.Username, invite.Code)
	log.Printf("No operators found. Use this link to create first operator: %s", inviteLink)

	return nil
}
