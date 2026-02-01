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

func (r *NewsRepository) CreateTx(tx *gorm.DB, news *model.News) error {
	return tx.Create(news).Error
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

func (r *NewsRepository) FindPaged(page, size int, status *int) ([]model.News, int64) {
	var news []model.News
	var total int64

	query := r.db.Model(&model.News{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Find(&news)

	return news, total
}

func (r *NewsRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
