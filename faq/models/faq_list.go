package models

type FAQList struct {
	Info         string        `yaml:"info"`
	InfoArticles []InfoArticle `yaml:"info_articles"`
	FAQs         []FAQ         `yaml:"faqs"`
}
