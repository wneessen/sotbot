package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	wftkaes "github.com/wneessen/go-wftk/crypto/aes"
	"github.com/wneessen/sotbot/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
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

func GetUsers(d *gorm.DB) ([]models.RegisteredUser, error) {
	userList := []models.RegisteredUser{}
	dbTx := d.Find(&userList)
	if dbTx.Error != nil {
		return userList, dbTx.Error
	}

	return userList, nil
}

func CreateUser(d *gorm.DB, u string) error {
	userObj := models.RegisteredUser{
		UserId: u,
	}
	dbTx := d.Create(&userObj)
	if dbTx.Error != nil {
		if dbTx.Error != nil {
			return dbTx.Error
		}
	}

	return nil
}

func DeleteUser(d *gorm.DB, u *models.RegisteredUser) error {
	dbTx := d.Unscoped().Delete(u)
	if dbTx.Error != nil {
		if dbTx.Error != nil {
			return dbTx.Error
		}
	}

	return nil
}

func UserSetPref(d *gorm.DB, u uint, k, v string) error {
	userPref := models.UserPref{}
	d.Where(models.UserPref{
		UserID: u,
		Key:    k,
	}).First(&userPref)

	if userPref.ID <= 0 {
		dbTx := d.Create(&models.UserPref{
			UserID: u,
			Key:    k,
			Value:  v,
		})
		if dbTx.Error != nil {
			return dbTx.Error
		}
		return nil
	}
	userPref.Value = v
	dbTx := d.Save(&userPref)
	if dbTx.Error != nil {
		return dbTx.Error
	}
	return nil
}

func UserSetPrefEnc(d *gorm.DB, c *viper.Viper, u uint, k, v string) error {
	encKey := c.GetString("enc_key")
	if encKey == "" {
		return fmt.Errorf("No encryption key set")
	}
	encVal, err := wftkaes.EncryptAuthBase64(v, encKey, k)
	if err != nil {
		return err
	}

	return UserSetPref(d, u, k, "{ENCRYPTED}"+encVal)
}

func UserDelPref(d *gorm.DB, u uint, k string) error {
	userPref := userGetPref(d, u, k)
	if userPref.ID <= 0 {
		return nil
	}

	dbTx := d.Unscoped().Delete(&userPref)
	if dbTx.Error != nil {
		return dbTx.Error
	}
	return nil
}

func UserGetPrefString(d *gorm.DB, u uint, k string) string {
	userPref := userGetPref(d, u, k)
	return userPref.Value
}

func UserGetPrefEncString(d *gorm.DB, c *viper.Viper, u uint, k string) string {
	l := log.WithFields(log.Fields{
		"action": "database.UserGetPrefEncString",
	})
	userPref := userGetPref(d, u, k)
	if !strings.HasPrefix(userPref.Value, "{ENCRYPTED}") {
		l.Warnf("Not an encrypted value")
		return userPref.Value
	}
	encKey := c.GetString("enc_key")
	if encKey == "" {
		l.Warnf("No encryption key set")
		return ""
	}
	decVal, err := wftkaes.DecryptAuthBase64(strings.Replace(userPref.Value, "{ENCRYPTED}", "", 1),
		encKey, k)
	if err != nil {
		l.Errorf("Failed to decrypt value: %v", err)
		return ""
	}

	return decVal
}

func userGetPref(d *gorm.DB, u uint, k string) models.UserPref {
	userPref := models.UserPref{}
	d.Where(models.UserPref{
		UserID: u,
		Key:    k,
	}).First(&userPref)

	return userPref
}
