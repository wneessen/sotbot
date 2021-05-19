package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func AnswerUser(s *discordgo.Session, m *discordgo.MessageCreate, msg string, mention string) {
	l := log.WithFields(log.Fields{
		"action": "bot.AnswerUser",
	})
	if mention != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, RandomArrr()+" "+mention+"! "+msg)
		if err != nil {
			l.Errorf("Failed to respond to author: %v", err)
		}
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, RandomArrr()+"! "+msg)
	if err != nil {
		l.Errorf("Failed to respond to author: %v", err)
	}

}
