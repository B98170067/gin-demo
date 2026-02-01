package repository

import (
	model "gin-demo/internal/models"

	"gorm.io/gorm"
)

type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{db}
}

func (r *NewsRepository) FindAll() ([]model.News, error) {
	var news []model.News
	err := r.db.Find(&news).Error
	return news, err
}

func (r *NewsRepository) Create(news *model.News) error {
	return r.db.Create(news).Error
}
