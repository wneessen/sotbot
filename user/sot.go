package user

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/database"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (u *User) UpdateSotBalance(d *gorm.DB, h *http.Client) error {
	l := log.WithFields(log.Fields{
		"action": "user.UpdateSotBalance",
	})

	l.Debugf("Checking if user %q has a RAT cookie set...", u.UserInfo.UserId)
	if !u.HasRatCookie() {
		l.Debugf("User %q has no cookie set.", u.UserInfo.UserId)
		return nil
	}

	userBalance, err := api.GetBalance(h, u.RatCookie)
	if err != nil {
		l.Errorf("Failed to fetch user balance from API: %v", err)
		return nil
	}

	if err := database.UpdateBalance(d, u.UserInfo.ID, &userBalance); err != nil {
		l.Errorf("Balance database update failed: %v", err)
	}

	return nil
}

func (u *User) UserNeedsNotifyFailedToken(d *gorm.DB) bool {
	l := log.WithFields(log.Fields{
		"action": "user.UserNeedsNotifyFailedToken",
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

func (u *User) CheckAuth(d *gorm.DB, h *http.Client) (bool, error) {
	l := log.WithFields(log.Fields{
		"action": "user.CheckAuth",
		"userId": u.UserInfo.ID,
	})
	l.Debugf("Checking SoT RAT cookie validtiy for user %v", u.UserInfo.UserId)

	failedRatTries, err := database.GetFailedRatCookieTries(d, u.UserInfo.ID)
	if err != nil {
		return false, err
	}
	if failedRatTries > 3 {
		l.Warnf("API requests with user's RAT cookie failed for more than 3 times. Skipping.")
		return false, nil
	}
	_, err = api.GetBalance(h, u.RatCookie)
	if err == nil {
		if err := database.UserDelPref(d, u.UserInfo.ID, "failed_rat_tries"); err != nil {
			l.Errorf("Failed to delete 'failed_rat_tries' userpref in DB: %v", err)
		}
		return false, nil
	}

	// If the response is a 403, it's likely that the RAT is expired or wrong.
	if err.Error() == "403" {
		newFailedTries, err := database.IncreaseFailedRatCookieTries(d, u.UserInfo.ID)
		if err != nil {
			return false, fmt.Errorf("Failed to increase rat_fails counter in DB: %v", err)
		}
		if newFailedTries > 3 {
			needsNotify := u.UserNeedsNotifyFailedToken(d)
			return needsNotify, nil
		}
	}

	return false, nil
}
