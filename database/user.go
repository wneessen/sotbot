package database

import (
	"fmt"
	"github.com/wneessen/sotbot/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUser(d *gorm.DB, u string) (models.RegisteredUser, error) {
	userObj := models.RegisteredUser{}
	dbTx := d.Preload(clause.Associations).Where(models.RegisteredUser{
		UserId: u,
	}).First(&userObj)
	if dbTx.Error != nil {
		if dbTx.Error.Error() != "record not found" {
			return userObj, fmt.Errorf("user lookup in DB failed: %v", dbTx.Error)
		}
	}

	return userObj, nil
}

func CreateUser(d *gorm.DB, u string, a bool) error {
	userObj := models.RegisteredUser{
		UserId:  u,
		IsAdmin: a,
	}
	dbTx := d.Create(&userObj)
	if dbTx.Error != nil {
		if dbTx.Error != nil {
			return dbTx.Error
		}
	}

	return nil
}
