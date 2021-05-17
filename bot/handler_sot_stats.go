package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/sotapi"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
		replyMsg := fmt.Sprintf("%v, thanks for setting/updating your RAT cookie.",
			m.Author.Mention())
		AnswerUser(s, m, replyMsg)
	}
}

// Get current SoT balance
func (b *Bot) GetBalance(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetBalance",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!balance") {
		l.Debugf("Received '!balance' request from user %v", m.Author.Username)
		userObj, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			l.Errorf("Database lookup failed: %v", err)
			return
		}
		if userObj.ID <= 0 {
			replyMsg := fmt.Sprintf("%v, sorry but your are not a registered user.",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
			return
		}
		userRatCookie := database.UserGetPrefString(b.Db, userObj.ID, "rat_cookie")
		if userRatCookie == "" {
			replyMsg := fmt.Sprintf("%v, sorry but you have no RAT cookie set. Try !setrat in the DMs",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
			return
		}

		userBalance, err := sotapi.GetBalance(b.HttpClient, userRatCookie)
		if err != nil {
			replyMsg := fmt.Sprintf("Sorry, %v but there was an error fetching your balance from the SoT API: %v",
				m.Author.Mention(), err)
			AnswerUser(s, m, replyMsg)
			return
		}

		p := message.NewPrinter(language.German)
		replyMsg := fmt.Sprintf("%v, your current SoT balance is: %v gold, %v doubloons and %v ancient coins",
			m.Author.Mention(),
			p.Sprintf("%d", userBalance.Gold),
			p.Sprintf("%d", userBalance.Doubloons),
			p.Sprintf("%d", userBalance.AncientCoins),
		)
		AnswerUser(s, m, replyMsg)
	}
}
