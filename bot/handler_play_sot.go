package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

func (b *Bot) UserPlaysSot(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	l := log.WithFields(log.Fields{
		"action": "handler.UserPlaySot",
	})

	reqUser, err := database.GetUser(b.Db, m.User.ID)
	if err != nil {
		l.Errorf("User lookup in DB failed: %v", err)
		return
	}
	if reqUser.ID <= 0 {
		l.Debugf("Received presence update, but user isn't registered.")
		return
	}
	discordUser, err := s.User(m.User.ID)
	if err != nil {
		l.Errorf("Discord user lookup failed: %v", err)
		return
	}

	// User started an activity, SoT maybe?
	if len(m.Activities) > 0 {
		for _, curActivity := range m.Activities {
			if curActivity.Name == "Sea of Thieves" {
				l.Debugf("%v started playing SoT. Updating balance...", discordUser.Username)
				b.UserUpdateSotBalance(&reqUser)

				if err := database.UserSetPref(b.Db, reqUser.ID, "playing_sot", time.Now().String()); err != nil {
					l.Errorf("Failed to update user status in database: %v", err)
				}
			}
		}
	}

	// User might have stopped an activity
	if len(m.Activities) == 0 {
		userWasPlaying := database.UserGetPrefString(b.Db, reqUser.ID, "playing_sot")
		if userWasPlaying != "" {
			l.Debugf("%v stopped playing SoT. Updating balance...", discordUser.Username)
			b.UserUpdateSotBalance(&reqUser)

			if err := database.UserDelPref(b.Db, reqUser.ID, "playing_sot"); err != nil {
				l.Errorf("Failed to delete user status in database: %v", err)
			}

			userBalance, err := database.GetBalance(b.Db, reqUser.ID)
			if err != nil {
				return
			}

			p := message.NewPrinter(language.German)
			dmText := fmt.Sprintf("%v, looks like you recently played SoT. Your new balance is: %v gold, %v "+
				"doubloons and %v ancient coins", discordUser.Mention(), p.Sprintf("%d", userBalance.Gold),
				p.Sprintf("%d", userBalance.Doubloons), p.Sprintf("%d", userBalance.AncientCoins))
			DmUser(s, reqUser.UserId, dmText)
		}
	}
}
