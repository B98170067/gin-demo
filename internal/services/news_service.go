package service

import (
	model "gin-demo/internal/models"
	repository "gin-demo/internal/repositories"
)

type NewsService struct {
	repo *repository.NewsRepository
}

func NewNewsService(repo *repository.NewsRepository) *NewsService {
	return &NewsService{repo}
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

func (s *NewsService) GetByID(id uint) (*model.News, error) {
	return s.repo.FindByID(id)
}

func (s *NewsService) Update(news *model.News) error {
	return s.repo.Update(news)
}

func (s *NewsService) Delete(id uint) error {
	return s.repo.Delete(id)
}
