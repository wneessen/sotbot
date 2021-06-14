package handler

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func TestCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.TestCmd",
	})
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	}); err != nil {
		l.Errorf("Failed: %v", err)
	}
}
