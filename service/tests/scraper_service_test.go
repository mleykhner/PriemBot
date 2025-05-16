package tests

import (
	"PriemBot/config"
	"PriemBot/scraper/mocks"
	scraperModels "PriemBot/scraper/models"
	"PriemBot/service"
	storageModels "PriemBot/storage/models"
	repoMocks "PriemBot/storage/repository/mocks"
	"errors"
	"testing"

	"gopkg.in/telebot.v4"
	"gorm.io/gorm"
)

func TestScraperService_TryRunAndNotify(t *testing.T) {
	tests := []struct {
		name            string
		mockSetup       func(*repoMocks.MockArticleRepository, *repoMocks.MockUserRepository, *mocks.MockBrowserlessClient)
		wantErr         bool
		wantNewArticles int
	}{
		{
			name: "successful scraping and notification",
			mockSetup: func(articleRepo *repoMocks.MockArticleRepository, userRepo *repoMocks.MockUserRepository, browserless *mocks.MockBrowserlessClient) {
				browserless.ExecuteScriptFunc = func(script string) ([]scraperModels.ScraperResultItem, error) {
					return []scraperModels.ScraperResultItem{
						{
							ID:       1,
							ImageURL: "http://example.com/image1.jpg",
							Title:    "Test Article 1",
							Subtitle: "Test Subtitle 1",
						},
						{
							ID:       2,
							ImageURL: "http://example.com/image2.jpg",
							Title:    "Test Article 2",
							Subtitle: "Test Subtitle 2",
						},
					}, nil
				}
				articleRepo.SaveAndGetNewFunc = func(results []storageModels.Article, tx *gorm.DB) ([]storageModels.Article, error) {
					return results, nil
				}
				userRepo.GetOperatorsFunc = func(tx *gorm.DB) ([]storageModels.User, error) {
					return []storageModels.User{
						{TelegramID: 123, Role: storageModels.RoleOperator},
						{TelegramID: 456, Role: storageModels.RoleOperator},
					}, nil
				}
			},
			wantErr:         false,
			wantNewArticles: 2,
		},
		{
			name: "browserless error",
			mockSetup: func(articleRepo *repoMocks.MockArticleRepository, userRepo *repoMocks.MockUserRepository, browserless *mocks.MockBrowserlessClient) {
				browserless.ExecuteScriptFunc = func(script string) ([]scraperModels.ScraperResultItem, error) {
					return nil, errors.New("browserless error")
				}
			},
			wantErr:         true,
			wantNewArticles: 0,
		},
		{
			name: "no new articles",
			mockSetup: func(articleRepo *repoMocks.MockArticleRepository, userRepo *repoMocks.MockUserRepository, browserless *mocks.MockBrowserlessClient) {
				browserless.ExecuteScriptFunc = func(script string) ([]scraperModels.ScraperResultItem, error) {
					return []scraperModels.ScraperResultItem{}, nil
				}
				articleRepo.SaveAndGetNewFunc = func(results []storageModels.Article, tx *gorm.DB) ([]storageModels.Article, error) {
					return []storageModels.Article{}, nil
				}
			},
			wantErr:         false,
			wantNewArticles: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticleRepo := &repoMocks.MockArticleRepository{}
			mockUserRepo := &repoMocks.MockUserRepository{}
			mockBrowserless := &mocks.MockBrowserlessClient{}
			tt.mockSetup(mockArticleRepo, mockUserRepo, mockBrowserless)
			mockConfig := &config.BrowserlessConfig{FilePath: ""}
			scraperService := service.NewScraperService(mockArticleRepo, mockUserRepo, &telebot.Bot{}, nil, mockConfig)
			scraperService.TryRunAndNotify()
		})
	}
}
