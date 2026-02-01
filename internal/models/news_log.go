package model

import "time"

type NewsLog struct {
	ID        uint `gorm:"primaryKey"`
	NewsID    uint
	Action    string
	CreatedAt time.Time
}
