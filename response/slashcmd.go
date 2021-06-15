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

func SlashCmdResponseEdit(s *discordgo.Session, i *discordgo.Interaction, u *user.User, msg string, mention bool) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdResponseDeferred",
	})
	if !mention {
		msg = fmt.Sprintf("%v! %v", RandomArrr(), msg)
	}
	if mention {
		msg = fmt.Sprintf("%v %v! %v", RandomArrr(), u.Mention, msg)
	}

	err := s.InteractionResponseEdit(s.State.User.ID, i, &discordgo.WebhookEdit{
		Content: msg,
	})
	if err != nil {
		l.Errorf("Failed respond to user's slash command request: %v", err)
	}
}

func SlashCmdEmbed(s *discordgo.Session, i *discordgo.Interaction, em *discordgo.MessageEmbed) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdEmbed",
	})

	emArr := make([]*discordgo.MessageEmbed, 0)
	emArr = append(emArr, em)
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: emArr,
		},
	})
	if err != nil {
		l.Errorf("Failed to send embed response to slash command: %v", err)
	}
}

func SlashCmdEmbedDeferred(s *discordgo.Session, i *discordgo.Interaction, em *discordgo.MessageEmbed) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdEmbed",
	})

	emArr := make([]*discordgo.MessageEmbed, 0)
	emArr = append(emArr, em)
	err := s.InteractionResponseEdit(s.State.User.ID, i, &discordgo.WebhookEdit{
		Embeds: emArr,
	})
	if err != nil {
		l.Errorf("Failed to send embed response to slash command: %v", err)
	}
}

func SlashCmdDel(s *discordgo.Session, i *discordgo.Interaction) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdDel",
	})
	err := s.InteractionResponseDelete(s.State.User.ID, i)
	if err != nil {
		l.Errorf("Failed to delete interaction response: %v", err)
	}
}
