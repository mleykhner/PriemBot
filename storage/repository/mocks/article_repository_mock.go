package mocks

import (
	"PriemBot/storage/models"

	"gorm.io/gorm"
)

type MockArticleRepository struct {
	SaveAndGetNewFunc func(results []models.Article, tx *gorm.DB) ([]models.Article, error)
}

func (m *MockArticleRepository) SaveAndGetNew(results []models.Article, tx *gorm.DB) ([]models.Article, error) {
	return m.SaveAndGetNewFunc(results, tx)
}
