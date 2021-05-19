package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"time"
)

// Let's the bot tell you the current date/time when requested via !time command
func (b *Bot) TellTime(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.TellTime",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.Content == "!time" {
		l.Debugf("Received '!time' request from user %v", m.Author.Username)
		returnMsg := fmt.Sprintf("The current time is: %v", time.Now().Format("2006-01-02 15:04:05 MST"))
		AnswerUser(s, m, returnMsg, m.Author.Mention())
	}
}
