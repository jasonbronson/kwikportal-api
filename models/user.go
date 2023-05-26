package models

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// User represents a user in the database.
type User struct {
	ID        string `gorm:"column:id"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate is a GORM callback that is triggered before creating a new user record.
// It generates a UUID for the ID field.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	u.ID = id.String()
	return nil
}

// TableName specifies the table name for the user model.
func (User) TableName() string {
	return "users"
}
