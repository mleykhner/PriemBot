package models

type FAQ struct {
	ID       int    `yaml:"id"`
	Question string `yaml:"q"`
	Answer   string `yaml:"a"`
}
