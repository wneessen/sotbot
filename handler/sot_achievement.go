package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
)

// Just a test handler
func GetSotAchievement(h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotAchievement",
	})

	userAchievement, err := api.GetLatestAchievement(h, u.RatCookie)
	if err != nil {
		l.Errorf("User's latest achievement lookup failed: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	embedTitle := fmt.Sprintf("%v, your latest achievement is: %v",
		u.AuthorName, userAchievement.Name)
	responseEmbed := discordgo.MessageEmbed{
		Title:       embedTitle,
		Description: userAchievement.Description,
		Image: &discordgo.MessageEmbedImage{
			URL: userAchievement.MediaUrl,
		},
		Type: discordgo.EmbedTypeImage,
	}
	return &responseEmbed, nil
}
