package bot

import (
	"PriemBot/service"
	"PriemBot/storage/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

type Handlers struct {
	bot            *TelegramBot
	userService    service.UserService
	dialogsService service.DialogsService
}

func NewBotHandlers(bot *TelegramBot, userService service.UserService, dialogsService service.DialogsService) *Handlers {
	return &Handlers{
		bot:            bot,
		userService:    userService,
		dialogsService: dialogsService,
	}
}

func (h *Handlers) RegisterHandlers() {
	bot := h.bot.GetBot()

	bot.Handle("/start", h.handleStartMessage)
	bot.Handle("/newop", h.handleNewOpMessage)
	// Обработка текстовых сообщений
	bot.Handle(tele.OnText, h.handleTextMessage)

	// Обработка callback-запросов (кнопки)
	bot.Handle(tele.OnCallback, h.handleCallback)
}

func (h *Handlers) handleNewOpMessage(c tele.Context) error {
	user, err := h.userService.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return err
	}

	if user.Role != models.RoleOperator {
		return c.Send("Я не знаю такую команду")
	}

	invite, err := h.userService.CreateInvite(c.Sender().ID)
	if err != nil {
		return err
	}

	link := h.userService.CreateInviteLink(invite.Code)
	message := fmt.Sprintf("Отправь эту ссылку новому оператору:\n%s\nОна будет действовать 24 часа", link)

	err = c.Send(message)

	return err
}

func (h *Handlers) handleStartMessage(c tele.Context) error {
	payload := c.Message().Payload
	user, err := h.userService.GetUserByTelegramID(c.Sender().ID)

	if err != nil {
		// Если пользователь не найден, создаем нового с ролью applicant
		user, err = h.userService.CreateUser(c.Sender().ID, c.Sender().Username, models.RoleApplicant)
		if err != nil {
			return c.Send("Произошла ошибка при регистрации пользователя")
		}
	}

	if len(payload) == 0 && user.Role == models.RoleApplicant {
		err := c.Send("Привет, абитуриент")
		if err != nil {
			return err
		}
	}

	if len(payload) == 0 && user.Role == models.RoleOperator {
		err := c.Send("Привет, оператор")
		if err != nil {
			return err
		}
	}

	if len(payload) > 0 && user.Role == models.RoleApplicant {
		err := h.userService.ApplyInvite(c.Sender().ID, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handlers) handleTextMessage(c tele.Context) error {
	// Получаем пользователя
	user, err := h.userService.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		// Если пользователь не найден, создаем нового с ролью applicant
		user, err = h.userService.CreateUser(c.Sender().ID, c.Sender().Username, models.RoleApplicant)
		if err != nil {
			return c.Send("Произошла ошибка при регистрации пользователя")
		}
	}

	// Если это команда /bye
	if c.Text() == "/bye" {
		return h.handleByeCommand(c, user)
	}

	// Если пользователь - абитуриент
	if user.Role == models.RoleApplicant {
		return h.handleApplicantMessage(c, user)
	}

	// Если пользователь - оператор
	if user.Role == models.RoleOperator {
		return h.handleOperatorMessage(c, user)
	}

	return nil
}

func (h *Handlers) handleApplicantMessage(c tele.Context, user *models.User) error {
	// Проверяем, есть ли активный диалог
	activeDialog, err := h.dialogsService.GetActiveDialogByStudentID(user.TelegramID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Send("Произошла ошибка при проверке диалога")
	}

	// Если нет активного диалога, создаем новый
	if activeDialog == nil {
		dialog, err := h.dialogsService.CreateDialog(user.TelegramID)
		if err != nil {
			return c.Send("Произошла ошибка при создании диалога")
		}

		// Получаем всех операторов
		operators, err := h.userService.GetOperators()
		if err != nil {
			return c.Send("Произошла ошибка при поиске операторов")
		}

		// Отправляем уведомления всем операторам
		for _, operator := range operators {
			msg, err := h.bot.GetBot().Send(
				&tele.User{ID: operator.TelegramID},
				fmt.Sprintf("Новый диалог от абитуриента @%s", user.Name),
				&tele.ReplyMarkup{
					InlineKeyboard: [][]tele.InlineButton{
						{
							{Text: "Принять", Data: fmt.Sprintf("accept_%d", dialog.ID)},
							{Text: "Отклонить", Data: fmt.Sprintf("decline_%d", dialog.ID)},
						},
					},
				},
			)
			if err != nil {
				continue
			}

			// Сохраняем уведомление в БД
			err = h.dialogsService.CreateDialogNotification(dialog.ID, msg.ID)
			if err != nil {
				continue
			}
		}

		return c.Send("Ваше сообщение отправлено операторам. Ожидайте ответа.")
	}

	// Если есть активный диалог, сохраняем сообщение
	_, err = h.dialogsService.CreateMessage(activeDialog.ID, user.TelegramID, c.Text())
	if err != nil {
		return c.Send("Произошла ошибка при сохранении сообщения")
	}

	// Отправляем сообщение оператору
	if activeDialog.OperatorID != nil {
		_, err := h.bot.GetBot().Send(
			&tele.User{ID: *activeDialog.OperatorID},
			fmt.Sprintf("Сообщение от @%s:\n%s", user.Name, c.Text()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handlers) handleOperatorMessage(c tele.Context, user *models.User) error {
	// Получаем активный диалог оператора
	activeDialog, err := h.dialogsService.GetActiveDialogByOperatorID(user.TelegramID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Send("Произошла ошибка при проверке диалога")
	}

	if activeDialog == nil {
		return c.Send("У вас нет активных диалогов")
	}

	// Сохраняем сообщение
	_, err = h.dialogsService.CreateMessage(activeDialog.ID, user.TelegramID, c.Text())
	if err != nil {
		return c.Send("Произошла ошибка при сохранении сообщения")
	}

	// Отправляем сообщение абитуриенту
	_, err = h.bot.GetBot().Send(
		&tele.User{ID: activeDialog.ApplicantID},
		c.Text(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handlers) handleCallback(c tele.Context) error {
	// Получаем пользователя
	user, err := h.userService.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Произошла ошибка при получении данных пользователя")
	}

	// Проверяем, что пользователь - оператор
	if user.Role != models.RoleOperator {
		return c.Send("Только операторы могут принимать диалоги")
	}

	// Парсим данные callback
	parts := strings.Split(c.Data(), "_")
	if len(parts) != 2 {
		return nil
	}

	action := parts[0]
	dialogID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return c.Send("Неверный формат данных")
	}

	if action == "accept" {
		// Назначаем оператора на диалог
		err = h.dialogsService.AssignOperator(uint(dialogID), user.TelegramID)
		if err != nil {
			return c.Send("Произошла ошибка при назначении оператора")
		}

		// Удаляем все уведомления о диалоге
		err = h.dialogsService.DeleteAllDialogNotifications(uint(dialogID))
		if err != nil {
			return c.Send("Произошла ошибка при удалении уведомлений")
		}

		// Отправляем сообщение оператору
		if err := c.Send("Вы приняли диалог. Теперь вы можете общаться с абитуриентом."); err != nil {
			return err
		}

		// Отправляем сообщение абитуриенту
		dialog, err := h.dialogsService.GetDialogByID(uint(dialogID))
		if err != nil {
			return err
		}

		_, err = h.bot.GetBot().Send(
			&tele.User{ID: dialog.ApplicantID},
			"Оператор принял ваш диалог. Теперь вы можете общаться.",
		)
		if err != nil {
			return err
		}

		err = c.Delete()

		if err != nil {
			return err
		}
	} else {
		err := c.Delete()
		if err != nil {
			return err
		}
	}

	// Отвечаем на callback
	return c.Respond()
}

func (h *Handlers) handleByeCommand(c tele.Context, user *models.User) error {
	var activeDialog *models.Dialog
	var err error

	if user.Role == models.RoleApplicant {
		activeDialog, err = h.dialogsService.GetActiveDialogByStudentID(user.TelegramID)
	} else {
		activeDialog, err = h.dialogsService.GetActiveDialogByOperatorID(user.TelegramID)
	}

	if err != nil {
		return c.Send("Произошла ошибка при получении диалога")
	}

	if activeDialog == nil {
		return c.Send("У вас нет активных диалогов")
	}

	// Закрываем диалог
	activeDialog.Status = models.DialogStatusClosed
	err = h.dialogsService.UpdateDialog(activeDialog)
	if err != nil {
		return c.Send("Произошла ошибка при закрытии диалога")
	}

	// Отправляем сообщение о закрытии диалога
	if err := c.Send("Диалог закрыт"); err != nil {
		return err
	}

	// Если это оператор, отправляем сообщение абитуриенту
	if user.Role == models.RoleOperator {
		_, err = h.bot.GetBot().Send(
			&tele.User{ID: activeDialog.ApplicantID},
			"Оператор закрыл диалог",
		)
		if err != nil {
			return err
		}
	} else {
		// Если это абитуриент, отправляем сообщение оператору
		if activeDialog.OperatorID != nil {
			_, err = h.bot.GetBot().Send(
				&tele.User{ID: *activeDialog.OperatorID},
				"Абитуриент закрыл диалог",
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
