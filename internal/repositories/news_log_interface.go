package repository

import (
	model "gin-demo/internal/models"

	"gorm.io/gorm"
)

//go:generate mockery --name=INewsLogRepository
type INewsLogRepository interface {
	CreateTx(tx *gorm.DB, log *model.NewsLog) error
}
