package models

import "gorm.io/gorm"

type Settings struct {
	gorm.Model
	UserID string `gorm:"column:user_id"`
}

func (Settings) TableName() string {
	return "settings"
}
