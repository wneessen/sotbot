package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/user"
	"gorm.io/gorm"
	"regexp"
)

// Self-check if a user is registered
func UserIsRegistered(u *user.User) string {
	if u.IsRegistered() {
		responeMsg := "You are a registered user!"
		return responeMsg
	}

	responeMsg := "You are not a registered user!"
	return responeMsg
}

// Register a new user
func RegisterUser(d *gorm.DB, u string) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.RegisterUser",
	})

	wrongFormatMsg := "Incorrect request format. Usage: !register <@user>"
	var validUser = regexp.MustCompile(`^<@[!&]*(\d+)>$`)
	if !validUser.MatchString(u) {
		return wrongFormatMsg, nil
	}
	validUserMatches := validUser.FindStringSubmatch(u)
	if len(validUserMatches) < 2 {
		return wrongFormatMsg, nil
	}
	dbUser, err := database.GetUser(d, validUserMatches[1])
	if err != nil {
		l.Errorf("Failed to look up user in database: %v", err)
		return "", fmt.Errorf("Could not look up if user exists in DB: %v", err)
	}
	if dbUser.ID > 0 {
		responseMsg := fmt.Sprintf("User %v is already registered.", validUserMatches[0])
		return responseMsg, nil
	}

	if err := database.CreateUser(d, validUserMatches[1]); err != nil {
		l.Errorf("Failed to store user in database: %v", err)
		return "", fmt.Errorf("Could not store user in DB: %v", err)
	}

	responseMsg := fmt.Sprintf("User %v successfully registered.", validUserMatches[0])
	return responseMsg, nil
}

// Unregister a user
func UnregisterUser(d *gorm.DB, u string) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.UnregisterUser",
	})

	wrongFormatMsg := "Incorrect request format. Usage: !unregister <@user>"
	var validUser = regexp.MustCompile(`^<@[!&]*(\d+)>$`)
	if !validUser.MatchString(u) {
		return wrongFormatMsg, nil
	}
	validUserMatches := validUser.FindStringSubmatch(u)
	if len(validUserMatches) < 2 {
		return wrongFormatMsg, nil
	}
	dbUser, err := database.GetUser(d, validUserMatches[1])
	if err != nil {
		l.Errorf("Failed to look up user in database: %v", err)
		return "", fmt.Errorf("Could not look up if user exists in DB: %v", err)
	}
	if dbUser.ID <= 0 {
		responseMsg := fmt.Sprintf("User %v is not registered.", validUserMatches[0])
		return responseMsg, nil
	}

	if err := database.DeleteUser(d, &dbUser); err != nil {
		l.Errorf("Failed to delete user in database: %v", err)
		return "", fmt.Errorf("Could not delete user in DB: %v", err)
	}

	responseMsg := fmt.Sprintf("User %v successfully unregistered.", validUserMatches[0])
	return responseMsg, nil
}
