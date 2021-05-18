package bot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/database/models"
	"time"
)

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

	st, err := b.Session.UserChannelCreate(u.UserId)
	if err != nil {
		l.Errorf("Failed to initiate DM channel with user: %v", err)
		return
	}
	du, err := b.Session.User(u.UserId)
	if err != nil {
		l.Errorf("Failed to look up discord user info: %v", err)
		return
	}
	replyMsg := fmt.Sprintf("%v, my last attempt to communicate with the SoT API failed. "+
		"This likely means, that your RAT cookie has expired. Please use the !setrat function to "+
		"update your cookie.", du.Mention())

	_, err = b.Session.ChannelMessageSend(st.ID, replyMsg)
	if err != nil {
		l.Errorf("Failed to notify user: %v", err)
		return
	}

	if err := database.UserSetPref(b.Db, u.ID, "failed_rat_notify", time.Now().String()); err != nil {
		l.Errorf("Failed to update user preference in DB: %v", err)
		return
	}
}
