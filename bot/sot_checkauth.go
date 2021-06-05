package bot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/random"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
	"time"
)

func (b *Bot) CheckSotAuth() {
	l := log.WithFields(log.Fields{
		"action": "bot.CheckSotAuth",
	})

	userList, err := database.GetUsers(b.Db)
	if err != nil {
		l.Errorf("Failed to fetch user list from DB: %v", err)
		return
	}
	for _, curUser := range userList {
		userObj, err := user.NewUser(b.Db, curUser.UserId)
		if err != nil {
			l.Errorf("Failed to create user object: %v", err)
			break
		}
		if userObj.HasRatCookie() {
			go func() {
				randNum, err := random.Number(600)
				if err != nil {
					l.Errorf("Failed to generate random number: %v", err)
					return
				}
				sleepTime, err := time.ParseDuration(fmt.Sprintf("%ds", randNum))
				if err != nil {
					l.Errorf("Failed to parse random number as duration: %v", err)
					return
				}
				time.Sleep(sleepTime)
				needsNotify, err := userObj.CheckAuth(b.Db, b.HttpClient)
				if err != nil {
					l.Errorf("CheckAuth failed: %v", err)
					return
				}
				if needsNotify {
					userObj.RatCookie = ""
					dmMsg := fmt.Sprintf("The last 3 attempts to communicate with the SoT API failed. " +
						"This likely means, that your RAT cookie has expired. Please use the !setrat function to " +
						"update your cookie.")
					response.DmUser(b.Session, userObj, dmMsg, true, false)
				}
			}()
		}
	}
}
