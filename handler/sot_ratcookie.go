package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/user"
	"gorm.io/gorm"
	"strconv"
)

type UserRat struct {
	Value      string `json:"Value"`
	Expiration int64  `json:"Expiration"`
}

// Set a SoT RAT cookie
func UserSetRatCookie(d *gorm.DB, c *viper.Viper, u *user.User, r string) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.UserSetRatCookie",
	})

	var ratCookieObj UserRat
	ratCookieBase64, err := base64.StdEncoding.DecodeString(r)
	steURL := "https://github.com/wneessen/sotbot-token-extrator"
	if err != nil {
		l.Errorf("Failed to decode base64: %v", err)
		return "", fmt.Errorf("Invalid input format. Please use the SoTBot Token Extractor at %v", steURL)
	}
	if err := json.Unmarshal([]byte(ratCookieBase64), &ratCookieObj); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return "", fmt.Errorf("Invalid input format. Please use the SoTBot Token Extractor at %v", steURL)
	}

	if err := database.UserSetPrefEnc(d, c, u.UserInfo.ID, "rat_cookie", ratCookieObj.Value); err != nil {
		l.Errorf("Failed to store RAT cookie in DB: %v", err)
		return "", err
	}
	expTimeString := strconv.FormatInt(ratCookieObj.Expiration, 10)
	if err := database.UserSetPrefEnc(d, c, u.UserInfo.ID, "rat_cookie_expire", expTimeString); err != nil {
		l.Errorf("Failed to store RAT cookie in DB: %v", err)
		return "", err
	}

	if err := database.UserDelPref(d, u.UserInfo.ID, "failed_rat_notify"); err != nil {
		l.Errorf("Failed to delete 'failed_rat_notify' preference: %v", err)
	}

	if err := database.UserDelPref(d, u.UserInfo.ID, "failed_rat_tries"); err != nil {
		l.Errorf("Failed to delete 'failed_rat_tries' userpref in DB: %v", err)
	}

	u.RatCookie = r
	responseMsg := "Thanks for setting/updating your RAT cookie."
	if err := database.UserDelPref(d, u.UserInfo.ID, "rat_expire_notify"); err != nil {
		l.Errorf("Failed to delete 'rat_expire_notify' user preference for user %q: %v", u.UserInfo.UserId,
			err)
	}
	return responseMsg, nil
}
