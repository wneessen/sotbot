package bot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/database/models"
	"github.com/wneessen/sotbot/sotapi"
	"time"
)

func (b *Bot) UpdateSotBalances() {
	l := log.WithFields(log.Fields{
		"action": "user.UpdateSotBalances",
	})
	l.Debugf("Looking for eligable users to update SoT balance...")

	userList, err := database.GetUsers(b.Db)
	if err != nil {
		l.Errorf("Failed to fetch registered users list: %v", err)
		return
	}

	for _, curUser := range userList {
		b.UserUpdateSotBalance(&curUser)
	}
}

func (b *Bot) UserUpdateSotBalance(u *models.RegisteredUser) {
	l := log.WithFields(log.Fields{
		"action": "user.UserUpdateSotBalance",
	})

	l.Debugf("Checking if user %q has a RAT cookie set...", u.UserId)
	userRatCookie := database.UserGetPrefString(b.Db, u.ID, "rat_cookie")
	if userRatCookie == "" {
		l.Debugf("User %q has not cookie set.", u.UserId)
		return
	}

	failedRatTries, err := database.GetFailedRatCookieTries(b.Db, u.ID)
	if err != nil {
		l.Errorf("Failed to fetch failed_rat_tries from DB: %v", err)
		return
	}
	if failedRatTries > 3 {
		l.Errorf("API requests with user's RAT cookie failed for more than 3 times. Skipping.")
		return
	}

	userBalance, err := sotapi.GetBalance(b.HttpClient, userRatCookie)
	if err != nil {
		if err.Error() == "403" {
			l.Errorf("Was not allowed to fetch user balance from API. Token likely invalid.")
			newFailedTries, err := database.IncreaseFailedRatCookieTries(b.Db, u.ID)
			if err != nil {
				l.Errorf("Failed to increase rat_fails counter in DB: %v", err)
				return
			}
			if newFailedTries > 3 {
				b.UserNotifyFailedToken(u)
			}

			return
		}
		l.Errorf("Failed to fetch user balance from API: %v", err)
		return
	}

	if err := database.UserDelPref(b.Db, u.ID, "failed_rat_tries"); err != nil {
		l.Errorf("Failed to delete 'failed_rat_tries' userpref in DB: %v", err)
	}

	if err := database.UpdateBalance(b.Db, u.ID, &userBalance); err != nil {
		l.Errorf("Balance database update failed: %v", err)
	}
}

func (b *Bot) UserNotifyFailedToken(u *models.RegisteredUser) {
	l := log.WithFields(log.Fields{
		"action": "user.UserNotifyFailedToken",
		"userId": u.UserId,
	})

	l.Debugf("Notifying user about broken token...")
	alreadyNotified := database.UserGetPrefString(b.Db, u.ID, "failed_rat_notify")
	if alreadyNotified != "" {
		l.Debugf("User %q has already been informed. Skipping.", u.UserId)
		return
	}

	dmText := fmt.Sprintf("The last 3 attempts to communicate with the SoT API failed. " +
		"This likely means, that your RAT cookie has expired. Please use the !setrat function to " +
		"update your cookie.")
	DmUser(b.Session, u.UserId, dmText, "")

	if err := database.UserSetPref(b.Db, u.ID, "failed_rat_notify", time.Now().String()); err != nil {
		l.Errorf("Failed to update user preference in DB: %v", err)
		return
	}
}