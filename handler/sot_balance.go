package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
	"net/http"
)

// Get current SoT balance
func GetSotBalance(d *gorm.DB, h *http.Client, u *user.User) (string, bool, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotBalance",
	})

	var retErr error
	retErr = u.UpdateSotBalance(d, h)
	userBalance, err := database.GetBalance(d, u.UserInfo.ID)
	if err != nil {
		l.Errorf("Database SoT balance lookup failed: %v", err)
		return "", false, err
	}

	p := message.NewPrinter(language.German)
	responseMsg := fmt.Sprintf("Your current SoT balance is: %v gold, %v doubloons and %v ancient coins",
		p.Sprintf("%d", userBalance.Gold), p.Sprintf("%d", userBalance.Doubloons),
		p.Sprintf("%d", userBalance.AncientCoins))
	return responseMsg, true, retErr
}
