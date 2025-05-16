package models

type FAQList struct {
	Info string `yaml:"info"`
	FAQs []FAQ  `yaml:"faqs"`
}
