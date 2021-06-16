package handler

import (
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"gorm.io/gorm"
)

// Register a new user
func RegisterUser(d *gorm.DB, u string) error {
	l := log.WithFields(log.Fields{
		"action": "handler.RegisterUser",
	})

	if err := database.CreateUser(d, u); err != nil {
		l.Errorf("Failed to store user in database: %v", err)
		return err
	}

	return nil
}

// Unregister a user
func UnregisterUser(d *gorm.DB, u string) error {
	l := log.WithFields(log.Fields{
		"action": "handler.UnregisterUser",
	})

	dbUser, err := database.GetUser(d, u)
	if err != nil {
		l.Errorf("Failed to look up user in database: %v", err)
		return err
	}
	if dbUser.ID <= 0 {
		return nil
	}
	if err := database.DeleteUser(d, &dbUser); err != nil {
		l.Errorf("Failed to delete user in database: %v", err)
		return err
	}

	return nil
}
