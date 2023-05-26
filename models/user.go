package models

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"column:id"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	u.ID = id.String()
	return nil
}

func (User) TableName() string {
	return "users"
}
