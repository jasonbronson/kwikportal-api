package repositories

import (
	"log"

	"github.com/jasonbronson/kwikportal-api/config"
	"github.com/jasonbronson/kwikportal-api/models"
)

func GetAllBookmarks() ([]models.Bookmark, error) {
	db := config.Cfg.GormDB

	var bookmarks []models.Bookmark
	result := db.Find(&bookmarks)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookmarks, nil
}

func GetUsersBookmarks(userID string) ([]models.Bookmark, error) {
	db := config.Cfg.GormDB

	var bookmarks []models.Bookmark
	result := db.Where("user_id=?", userID).Find(&bookmarks)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookmarks, nil
}

func SaveAllBookmarks(bookmarks []models.Bookmark) error {
	db := config.Cfg.GormDB

	result := db.Table("bookmarks").Create(&bookmarks)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return result.Error
	}

	return nil
}
