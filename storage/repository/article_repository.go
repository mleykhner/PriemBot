package repository

import (
	"PriemBot/storage/models"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	SaveAndGetNew(results []models.Article, tx *gorm.DB) ([]models.Article, error)
}

type ArticleRepositoryImpl struct{}

func NewArticleRepository() ArticleRepository {
	return &ArticleRepositoryImpl{}
}

func (a ArticleRepositoryImpl) SaveAndGetNew(results []models.Article, tx *gorm.DB) ([]models.Article, error) {
	if len(results) == 0 {
		return nil, nil
	}

	ids := make([]uint, 0, len(results))
	for _, art := range results {
		ids = append(ids, art.ID)
	}

	var existingIDs []uint
	if err := tx.Model(&models.Article{}).
		Where("id IN ?", ids).
		Pluck("id", &existingIDs).Error; err != nil {
		return nil, err
	}

	existingSet := make(map[uint]struct{}, len(existingIDs))
	for _, id := range existingIDs {
		existingSet[id] = struct{}{}
	}

	var newArticles []models.Article
	for _, art := range results {
		if _, exists := existingSet[art.ID]; !exists {
			newArticles = append(newArticles, art)
		}
	}

	if len(newArticles) > 0 {
		if err := tx.Create(&newArticles).Error; err != nil {
			return nil, err
		}
	}

	return newArticles, nil
}
