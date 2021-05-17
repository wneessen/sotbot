package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func AnswerUser(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	l := log.WithFields(log.Fields{
		"action": "bot.AnswerUser",
	})
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		l.Errorf("Failed to respond to author: %v", err)
	}
}
