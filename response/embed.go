package response

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func Embed(s *discordgo.Session, c string, em *discordgo.MessageEmbed) {
	l := log.WithFields(log.Fields{
		"action": "bot.Embed",
	})

	_, err := s.ChannelMessageSendEmbed(c, em)
	if err != nil {
		l.Errorf("Failed to send embeded message: %v", err)
	}
}
