package service

import (
	model "gin-demo/internal/models"
	repository "gin-demo/internal/repositories"

	"gorm.io/gorm"
)

type NewsService struct {
	db      *gorm.DB
	repo    *repository.NewsRepository
	logRepo *repository.NewsLogRepository
}

func NewNewsService(
	db *gorm.DB,
	repo *repository.NewsRepository,
	logRepo *repository.NewsLogRepository,
) *NewsService {
	return &NewsService{
		db:      db,
		repo:    repo,
		logRepo: logRepo,
	}
}

func (s *NewsService) GetAllNews() ([]model.News, error) {
	return s.repo.FindAll()
}

func (s *NewsService) GetPaged(page, size int, status *int) ([]model.News, int64) {
	return s.repo.FindPaged(page, size, status)
}

func (s *NewsService) CreateNews(news *model.News) error {
	return s.repo.Create(news)
}

func (s *NewsService) CreateWithLog(news *model.News) error {
	return s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateTx(tx, news); err != nil {
			return err // rollback
		}

		log := &model.NewsLog{
			NewsID: news.ID,
			Action: "CREATE",
		}

		if err := s.logRepo.CreateTx(tx, log); err != nil {
			return err // rollback
		}

		return nil // commit
	})
}

func (s *NewsService) GetByID(id uint) (*model.News, error) {
	return s.repo.FindByID(id)
}

func (s *NewsService) Update(news *model.News) error {
	return s.repo.Update(news)
}

func (s *NewsService) Delete(id uint) error {
	return s.repo.Delete(id)
}
