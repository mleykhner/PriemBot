package service

import (
	"PriemBot/scraper/models"
	models2 "PriemBot/storage/models"
	"PriemBot/storage/repository"
	"bytes"
	"encoding/json"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"time"
)

type ScraperService struct {
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
	bot         *tele.Bot
	db          *gorm.DB
	jsFilePath  string
}

func NewScraperService(articleRepo repository.ArticleRepository, userRepo repository.UserRepository, bot *tele.Bot, db *gorm.DB, jsFilePath string) *ScraperService {
	return &ScraperService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
		bot:         bot,
		db:          db,
		jsFilePath:  jsFilePath,
	}
}

func (s *ScraperService) StartScheduler() {
	// Первый запуск сразу
	go s.TryRunAndNotify()
	// Два раза в сутки
	ticker := time.NewTicker(12 * time.Hour)
	go func() {
		for range ticker.C {
			s.TryRunAndNotify()
		}
	}()
}

func (s *ScraperService) TryRunAndNotify() {
	js, err := os.ReadFile(s.jsFilePath)
	if err != nil {
		return
	}
	results, err := RunBrowserless(string(js)) // функцию см. выше
	if err != nil {
		return
	}
	var articles []models2.Article
	for _, result := range results {
		articles = append(articles, result.ToArticle())
	}
	newItems, err := s.articleRepo.SaveAndGetNew(articles, s.db)
	if err != nil || len(newItems) == 0 {
		return
	}
	// Получить список всех юзеров:
	users, _ := s.userRepo.GetApplicants(s.db) // или другой сервис Users
	for _, u := range users {
		for _, item := range newItems {
			text := fmt.Sprintf("*%s*\n\n%s\n\n[Читать подробнее](https://mai.ru/press/news/detail.php?ID=%d)", item.Title, item.Subtitle, item.ID)
			photo := &tele.Photo{File: tele.FromURL(item.ImageURL), Caption: text}
			_, err := s.bot.Send(&tele.User{ID: u.TelegramID}, photo, tele.ModeMarkdown)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func RunBrowserless(js string) ([]models.ScraperResultItem, error) {
	reqBody, _ := json.Marshal(models.Request{Code: string(js)})
	resp, err := http.Post("http://localhost:3000/function", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var result []models.ScraperResultItem
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
