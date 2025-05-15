package bot

import (
	"PriemBot/config"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v4"
)

type TelegramBot struct {
	bot *tele.Bot
}

func NewTelegramBot(config *config.BotConfig) (*TelegramBot, error) {
	if config.Token == "" {
		return nil, fmt.Errorf("bot token is required")
	}

	pref := tele.Settings{
		Token:  config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &TelegramBot{bot: bot}, nil
}

func (t *TelegramBot) GetBot() *tele.Bot {
	return t.bot
}

func (t *TelegramBot) Start() {
	t.bot.Start()
}

func (t *TelegramBot) Stop() {
	t.bot.Stop()
}
