package bot

import (
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
	"time"
)

func (b *Bot) CheckRatCookies() {
	l := log.WithFields(log.Fields{
		"action": "bot.CheckRatCookies",
	})

	userList, err := database.GetUsers(b.Db)
	if err != nil {
		l.Errorf("Failed to fetch user list from DB: %v", err)
		return
	}
	for _, curUser := range userList {
		userObj, err := user.NewUser(b.Db, b.Config, curUser.UserId)
		if err != nil {
			l.Errorf("Failed to create user object: %v", err)
			continue
		}
		if userObj.HasRatCookie() && !userObj.RatIsValid() {
			userNotified := database.UserGetPrefString(b.Db, userObj.UserInfo.ID, "rat_expire_notify")
			if userNotified == "" {
				dmMsg := "Your SoT RAT cookie has expired. Please use the `/setrat` command to set a new one."
				response.DmUser(b.Session, userObj, dmMsg, true, false)
				if err := database.UserSetPref(b.Db, userObj.UserInfo.ID, "rat_expire_notify",
					time.Now().String()); err != nil {
					l.Errorf("Failed to set 'rat_expire_notify' flag in database for user %q: %v",
						userObj.UserInfo.UserId, err)
				}
			}
		}
	}
}
