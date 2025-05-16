package models

import "PriemBot/storage/models"

type ScraperResultItem struct {
	ID       int    `json:"id,string"`
	ImageURL string `json:"imageUrl"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

func (s ScraperResultItem) ToArticle() models.Article {
	return models.Article{
		ID:       uint(s.ID),
		ImageURL: s.ImageURL,
		Title:    s.Title,
		Subtitle: s.Subtitle,
	}
}
