package user

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/database/models"
	"gorm.io/gorm"
)

type User struct {
	AuthorId       string
	AuthorName     string
	ChanPermission int64
	Mention        string
	RatCookie      string
	UserInfo       *models.RegisteredUser
}

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

func (u *User) IsRegistered() bool {
	return u.UserInfo.ID > 0
}

func (u *User) HasRatCookie() bool {
	return u.RatCookie != ""
}

func (u *User) IsAdmin() bool {
	return u.ChanPermission&discordgo.PermissionAdministrator != 0
}
