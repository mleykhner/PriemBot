package config

type BrowserlessConfig struct {
	FilePath string
	Host     string
}

func NewBrowserlessConfig() *BrowserlessConfig {
	return &BrowserlessConfig{
		FilePath: getEnvOrDefault("BROWSERLESS_FILE_PATH", ""),
		Host:     getEnvOrDefault("BROWSERLESS_HOST", ""),
	}
}
