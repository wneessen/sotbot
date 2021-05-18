package bot

import (
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/sotapi"
	"time"
)

func (b *Bot) UpdateSotBalance(t time.Time) {
	l := log.WithFields(log.Fields{
		"action": "timer.UpdateSotBalance",
	})
	l.Debugf("Looking for eligable users to update SoT balance...")

	userList, err := database.GetUsers(b.Db)
	if err != nil {
		l.Errorf("Failed to fetch registered users list: %v", err)
		return
	}

	for _, curUser := range userList {
		l.Debugf("Checking if user %q has a RAT cookie set...", curUser.UserId)
		userRatCookie := database.UserGetPrefString(b.Db, curUser.ID, "rat_cookie")
		if userRatCookie == "" {
			l.Debugf("User %q has not cookie set.", curUser.UserId)
			break
		}
		brokenCookie := database.UserGetPrefString(b.Db, curUser.ID, "failed_rat_notify")
		if brokenCookie != "" {
			l.Debugf("User's RAT cookie was broken at last attempt. Skipping.")
			break
		}

		userBalance, err := sotapi.GetBalance(b.HttpClient, userRatCookie)
		if err != nil {
			if err.Error() == "403" {
				l.Errorf("Was not allowed to fetch user balance from API. Token likely invalid.")
				b.UserNotifyFailedToken(&curUser)
				return
			}
			l.Errorf("Failed to fetch user balance from API: %v", err)
			return
		}

		if err := database.UpdateBalance(b.Db, curUser.ID, &userBalance); err != nil {
			l.Errorf("Balance database update failed: %v", err)
			break
		}
	}
}
