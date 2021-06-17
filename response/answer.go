package response

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func AnswerUser(s *discordgo.Session, m *discordgo.MessageCreate, msg string, mention bool) {
	l := log.WithFields(log.Fields{
		"action": "response.AnswerUser",
	})
	if mention {
		_, err := s.ChannelMessageSend(m.ChannelID, RandomArrr()+" "+m.Author.Mention()+"! "+msg)
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
