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
		returnMsg := fmt.Sprintf("%v, the current time is: %v",
			m.Author.Mention(), time.Now().Format("2006-01-02 15:04:05 MST"))
		_, err := s.ChannelMessageSend(m.ChannelID, returnMsg)
		if err != nil {
			l.Errorf("Failed to respond to author: %v", err)
		}
	}
}
