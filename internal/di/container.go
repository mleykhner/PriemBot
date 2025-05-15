package di

import (
	"PriemBot/config"
	"PriemBot/infrastructure/bot"
	"PriemBot/infrastructure/database"
	"PriemBot/service"
	"PriemBot/storage/repository"
)

type Container struct {
	UserService    service.UserService
	DialogsService service.DialogsService
	InitService    service.InitService
	DB             *database.PostgresDB
	Bot            *bot.TelegramBot
	BotHandlers    *bot.Handlers
}

func NewContainer() (*Container, error) {
	// Инициализация конфигурации
	dbConfig := config.NewDatabaseConfig()
	botConfig := config.NewBotConfig()

	// Инициализация базы данных
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		return nil, err
	}

	// Выполнение миграций
	if err := db.AutoMigrate(); err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	inviteRepo := repository.NewInviteRepository()

	// Инициализация сервисов
	userService := service.NewUserService(db.GetDB(), botConfig)
	dialogsService := service.NewDialogsService(db.GetDB())
	initService := service.NewInitService(userService, inviteRepo, db.GetDB())

	// Проверка наличия операторов и создание приглашения при необходимости
	if err := initService.CheckAndCreateOperator(botConfig); err != nil {
		return nil, err
	}

	// Инициализация бота
	telegramBot, err := bot.NewTelegramBot(botConfig)
	if err != nil {
		return nil, err
	}

	// Инициализация обработчиков бота
	botHandlers := bot.NewBotHandlers(telegramBot, userService, dialogsService)
	botHandlers.RegisterHandlers()

	return &Container{
		UserService:    userService,
		DialogsService: dialogsService,
		InitService:    initService,
		DB:             db,
		Bot:            telegramBot,
		BotHandlers:    botHandlers,
	}, nil
}

func (c *Container) Close() error {
	// Останавливаем бота
	c.Bot.Stop()

	// Закрываем соединение с базой данных
	return c.DB.Close()
}
