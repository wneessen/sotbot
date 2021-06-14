package response

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/user"
)

func SlashCmdResponse(s *discordgo.Session, i *discordgo.Interaction, u *user.User, msg string, mention bool) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdResponse",
	})
	if !mention {
		msg = fmt.Sprintf("%v! %v", RandomArrr(), msg)
	}
	if mention {
		msg = fmt.Sprintf("%v %v! %v", RandomArrr(), u.Mention, msg)
	}

	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: msg,
		},
	})
	if err != nil {
		l.Errorf("Failed respond to user's slash command request: %v", err)
	}
}

func SlashCmdResponseDeferred(s *discordgo.Session, i *discordgo.Interaction) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdResponseDeferred",
	})

	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "",
		},
	})
	if err != nil {
		l.Errorf("Failed respond to user's slash command request: %v", err)
	}
}
