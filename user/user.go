// Package user provides a User object and it's associated methods
package user

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/database/models"
	"gorm.io/gorm"
)

// User is a struct of user object. Each incoming message in the command_handler
// will try to create a User object based on the message it received
type User struct {
	// AuthorId is the Discord UserId of the message current session event
	AuthorId string

	// AuthorName is the Discord nickname of the AuthorId
	AuthorName string

	// ChanPermission reflects the bitmask of permissions of User in the
	// current channel
	ChanPermission int64

	// Mention is the Discord string for a @mention of the User
	Mention string

	// RatCookie reflects the SoT RAT authentication cookie from the database
	// if the user set such before
	RatCookie string

	// UserInfo reflects a pointer to the *models.RegisteredUser database model of
	// User
	UserInfo *models.RegisteredUser
}

// NewUser creates and returns a new User object
func NewUser(d *gorm.DB, c *viper.Viper, i string) (*User, error) {
	l := log.WithFields(log.Fields{
		"action": "user.NewUser",
	})
	dbUser, err := database.GetUser(d, i)
	if err != nil {
		l.Errorf("Database user lookup failed: %v", err)
		return &User{}, err
	}

	userObj := User{AuthorId: i, UserInfo: &dbUser}

	if dbUser.ID > 0 {
		userRatCookie := database.UserGetPrefEncString(d, c, dbUser.ID, "rat_cookie")
		userObj.RatCookie = userRatCookie
	}

	return &userObj, nil
}

// IsRegistered checks wether *User is registered or not
func (u *User) IsRegistered() bool {
	return u.UserInfo.ID > 0
}

// HasRatCookie checks wether *User has a valid SoT RAT cookie set in the database
func (u *User) HasRatCookie() bool {
	return u.RatCookie != ""
}

// IsAdmin checks weather *User has channel admin permission
func (u *User) IsAdmin() bool {
	return u.ChanPermission&discordgo.PermissionAdministrator != 0
}
