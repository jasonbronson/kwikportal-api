package models

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Bookmark represents a bookmark entry in the database.
type Bookmark struct {
	ID        string `gorm:"column:id"`
	UserID    string `gorm:"column:user_id"`
	Folder    string `gorm:"column:folder"`
	URL       string `gorm:"column:url"`
	AddDate   int64  `gorm:"column:add_date"`
	Icon      string `gorm:"column:icon"`
	Name      string `gorm:"column:name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate is a GORM callback that is triggered before creating a new bookmark record.
// It generates a UUID for the ID field.
func (b *Bookmark) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	b.ID = id.String()
	return nil
}

// TableName specifies the table name for the bookmark model.
func (Bookmark) TableName() string {
	return "bookmarks"
}
