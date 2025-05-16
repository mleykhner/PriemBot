package config

type FAQConfig struct {
	FilePath string
}

func NewFAQConfig() *FAQConfig {
	return &FAQConfig{
		FilePath: getEnvOrDefault("FAQ_FILE_PATH", ""),
	}
}
