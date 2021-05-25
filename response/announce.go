package response

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func Announce(s *discordgo.Session, c, m string) {
	l := log.WithFields(log.Fields{
		"action": "bot.Announce",
	})
	_, err := s.ChannelMessageSend(c, RandomArrr()+"! "+m)
	if err != nil {
		l.Errorf("Failed to make announcement: %v", err)
	}
}
