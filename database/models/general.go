package models

import "gorm.io/gorm"

type General struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt int64          `gorm:"autoCreateTime"`
	UpdatedAt int64          `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
