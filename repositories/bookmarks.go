package repositories

import (
	"log"

	"github.com/jasonbronson/kwikportal-api/config"
	"github.com/jasonbronson/kwikportal-api/models"
)

// GetAllBookmarks retrieves all bookmarks from the database.
func GetAllBookmarks() ([]models.Bookmark, error) {
	db := config.Cfg.GormDB

	var bookmarks []models.Bookmark
	result := db.Find(&bookmarks)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookmarks, nil
}

// GetUsersBookmarks retrieves bookmarks associated with a specific user.
func GetUsersBookmarks(userID string) ([]models.Bookmark, error) {
	db := config.Cfg.GormDB

	var bookmarks []models.Bookmark
	result := db.Where("user_id=?", userID).Find(&bookmarks)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookmarks, nil
}

// SaveAllBookmarks saves multiple bookmarks to the database.
func SaveAllBookmarks(bookmarks []models.Bookmark) error {
	db := config.Cfg.GormDB

	result := db.Debug().Table("bookmarks").Create(&bookmarks)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return result.Error
	}

	return nil
}

// SaveBookmark saves a single bookmark row to the database
func SaveBookmark(bookmark models.Bookmark, userID string) error {
	db := config.Cfg.GormDB

	result := db.Debug().Table("bookmarks").Where("id = ? AND user_id = ?", bookmark.ID, userID).Updates(&bookmark)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return result.Error
	}

	return nil
}

func DeleteBookmark(bookmarkID string, userID string) error {
	db := config.Cfg.GormDB

	// Delete the bookmark based on the bookmark ID and user ID
	result := db.Debug().Table("bookmarks").Where("id = ? AND user_id = ?", bookmarkID, userID).Delete(&models.Bookmark{})
	if result.Error != nil {
		log.Println(result.Error.Error())
		return result.Error
	}

	return nil
}
