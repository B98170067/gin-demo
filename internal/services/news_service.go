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

func (s *NewsService) CreateNews(news *model.News) error {
	return s.repo.Create(news)
}
