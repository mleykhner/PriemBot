package config

type BotConfig struct {
	Token    string
	Username string
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		Token:    getEnvOrDefault("BOT_TOKEN", ""),
		Username: getEnvOrDefault("BOT_USERNAME", ""),
	}
}
