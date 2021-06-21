// Package user provides a User object and it's associated methods
package user

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/database/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// User is a struct of user object. Each incoming message in the command_handler
// will try to create a User object based on the message it received
type User struct {
	AuthorId       string
	AuthorName     string
	ChanPermission int64
	Mention        string
	RatCookie      string
	RatExpire      time.Time
	UserInfo       *models.RegisteredUser
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

		userRatCookieExpireString := database.UserGetPrefEncString(d, c, dbUser.ID, "rat_cookie_expire")
		userRatCookieExpireInt, err := strconv.ParseInt(userRatCookieExpireString, 10, 64)
		if err != nil {
			l.Errorf("Failed to convert string to int64: %v", err)
		}
		userObj.RatExpire = time.Unix(userRatCookieExpireInt, 0)
	}

	return &userObj, nil
}

// IsRegistered return true when *User is registered
func (u *User) IsRegistered() bool {
	return u.UserInfo.ID > 0
}

// HasRatCookie returns true when *User has a valid SoT RAT cookie set in
// the database
func (u *User) HasRatCookie() bool {
	return u.RatCookie != ""
}

// RatIsValid returns true when the RAT cookie of *User is still valid
func (u *User) RatIsValid() bool {
	return u.RatExpire.Unix() > time.Now().Unix()
}

// IsAdmin returns true when *User has channel admin permission
func (u *User) IsAdmin() bool {
	return u.ChanPermission&discordgo.PermissionAdministrator != 0
}
