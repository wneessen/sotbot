package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"strings"
)

// Set a SoT RAT cookie
func (b *Bot) SetRatCookie(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.SetRatCookie",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	curChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		l.Errorf("Failed to get channel info: %v", err)
		return
	}
	if curChannel.Type != discordgo.ChannelTypeDM {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!setrat") {
		l.Debugf("Received '!setrat' request from user %v", m.Author.Username)
		userObj, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			l.Errorf("Database lookup failed: %v", err)
			return
		}

		var returnMsg string
		if userObj.ID <= 0 {
			returnMsg = "Sorry, you are not a registered user."
			AnswerUser(s, m, returnMsg)
			return
		}

		wrongFormatMsg := fmt.Sprintf("%v, incorrect request format. Usage: !setrat <ratcookie>",
			m.Author.Mention())
		msgArray := strings.SplitN(m.Message.Content, " ", 2)
		if len(msgArray) != 2 {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}

		if err := database.UserSetPref(b.Db, userObj.ID, "rat_cookie", msgArray[1]); err != nil {
			l.Errorf("Failed to store RAT cookie in DB: %v", err)
			replyMsg := fmt.Sprintf("Sorry, I couldn't store/update your cookie in the DB.")
			AnswerUser(s, m, replyMsg)
			return
		}

		if err := database.UserDelPref(b.Db, userObj.ID, "failed_rat_notify"); err != nil {
			l.Errorf("Failed to delete 'failed_rat_notify' preference: %v", err)
		}

		replyMsg := fmt.Sprintf("%v, thanks for setting/updating your RAT cookie.",
			m.Author.Mention())
		AnswerUser(s, m, replyMsg)
	}
}
