package user

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/sotapi"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (u *User) UpdateSotBalance(d *gorm.DB, h *http.Client) error {
	l := log.WithFields(log.Fields{
		"action": "user.UpdateSotBalance",
	})

	l.Debugf("Checking if user %q has a RAT cookie set...", u.UserInfo.UserId)
	userRatCookie := database.UserGetPrefString(d, u.UserInfo.ID, "rat_cookie")
	if userRatCookie == "" {
		l.Debugf("User %q has not cookie set.", u.UserInfo.UserId)
		return nil
	}

	failedRatTries, err := database.GetFailedRatCookieTries(d, u.UserInfo.ID)
	if err != nil {
		l.Errorf("Failed to fetch failed_rat_tries from DB: %v", err)
		return nil
	}
	if failedRatTries > 3 {
		l.Errorf("API requests with user's RAT cookie failed for more than 3 times. Skipping.")
		return nil
	}

	userBalance, err := sotapi.GetBalance(h, userRatCookie)
	if err != nil {
		if err.Error() == "403" {
			l.Errorf("Was not allowed to fetch user balance from API. Token likely invalid.")
			newFailedTries, err := database.IncreaseFailedRatCookieTries(d, u.UserInfo.ID)
			if err != nil {
				l.Errorf("Failed to increase rat_fails counter in DB: %v", err)
				return nil
			}
			if newFailedTries > 3 {
				needsNotify := u.UserNotifyFailedToken(d)
				if needsNotify {
					return fmt.Errorf("notify")
				}
			}
			return nil
		}
		l.Errorf("Failed to fetch user balance from API: %v", err)
		return nil
	}

	if err := database.UserDelPref(d, u.UserInfo.ID, "failed_rat_tries"); err != nil {
		l.Errorf("Failed to delete 'failed_rat_tries' userpref in DB: %v", err)
	}

	if err := database.UpdateBalance(d, u.UserInfo.ID, &userBalance); err != nil {
		l.Errorf("Balance database update failed: %v", err)
	}

	return nil
}

func (u *User) UserNotifyFailedToken(d *gorm.DB) bool {
	l := log.WithFields(log.Fields{
		"action": "user.UserNotifyFailedToken",
		"userId": u.UserInfo.ID,
	})

	l.Debugf("Notifying user about broken token...")
	alreadyNotified := database.UserGetPrefString(d, u.UserInfo.ID, "failed_rat_notify")
	if alreadyNotified != "" {
		l.Debugf("User %q has already been informed. Skipping.", u.UserInfo.UserId)
		return false
	}

	if err := database.UserSetPref(d, u.UserInfo.ID, "failed_rat_notify", time.Now().String()); err != nil {
		l.Errorf("Failed to update user preference in DB: %v", err)
		return false
	}

	return true
}
