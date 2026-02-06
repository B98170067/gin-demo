package repository

import (
	model "gin-demo/internal/models"

	"gorm.io/gorm"
)

//go:generate mockery --name=INewsRepository
type INewsRepository interface {
	FindAll() ([]model.News, error)
	FindPaged(page, size int, status *int) ([]model.News, int64)
	Create(news *model.News) error
	CreateTx(tx *gorm.DB, news *model.News) error
	FindByID(id uint) (*model.News, error)
	Update(news *model.News) error
	Delete(id uint) error
	Transaction(fn func(tx *gorm.DB) error) error
}
