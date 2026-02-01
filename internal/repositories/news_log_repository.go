package repository

import (
	model "gin-demo/internal/models"

	"gorm.io/gorm"
)

type NewsLogRepository struct {
	db *gorm.DB
}

func NewNewsLogRepository(db *gorm.DB) *NewsLogRepository {
	return &NewsLogRepository{db}
}

func (r *NewsLogRepository) CreateTx(tx *gorm.DB, log *model.NewsLog) error {
	return tx.Create(log).Error
}
