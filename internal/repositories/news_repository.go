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

func (r *NewsRepository) FindByID(id uint) (*model.News, error) {
	var news model.News
	err := r.db.First(&news, id).Error
	return &news, err
}

func (r *NewsRepository) Update(news *model.News) error {
	return r.db.Save(news).Error
}

func (r *NewsRepository) Delete(id uint) error {
	return r.db.Delete(&model.News{}, id).Error
}
