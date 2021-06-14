package response

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func SlashCmdEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, em *discordgo.MessageEmbed) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdEmbed",
	})

	emArr := make([]*discordgo.MessageEmbed, 0)
	emArr = append(emArr, em)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: emArr,
		},
	})
	if err != nil {
		l.Errorf("Failed to send embed response to slash command: %v", err)
	}
}

func SlashCmdEmbedDeferred(s *discordgo.Session, i *discordgo.InteractionCreate, em *discordgo.MessageEmbed) {
	l := log.WithFields(log.Fields{
		"action": "response.SlashCmdEmbed",
	})

	emArr := make([]*discordgo.MessageEmbed, 0)
	emArr = append(emArr, em)
	err := s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
		Embeds: emArr,
	})
	if err != nil {
		l.Errorf("Failed to send embed response to slash command: %v", err)
	}
}
