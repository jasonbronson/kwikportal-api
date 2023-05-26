package repositories

import (
	"github.com/jasonbronson/kwikportal-api/config"
	"github.com/jasonbronson/kwikportal-api/models"
)

func GetUser(email string) (models.User, error) {
	db := config.Cfg.GormDB

	var foundUser models.User
	result := db.Debug().Where("email = ?", email).First(&foundUser)
	if result.Error != nil {
		return foundUser, result.Error
	}

	return foundUser, nil
}
func SaveUser(user models.User) error {
	db := config.Cfg.GormDB

	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
