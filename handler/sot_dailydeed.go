package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
)

// Get the daily deed
func GetDailyDeed(h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetDailyDeed",
	})
	dailyDeed, err := api.GetDailyDeed(h, u.RatCookie)
	if err != nil {
		l.Errorf("An error occured fetching daily deeed data: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	embedDesc := fmt.Sprintf("%v\n\nStart date: %v\nEnd date: %v", dailyDeed.Copy,
		dailyDeed.StartDate.Time().Format("2006-01-02 15:04:05 MST"),
		dailyDeed.EndDate.Time().Format("2006-01-02 15:04:05 MST"))
	responseEmbed := discordgo.MessageEmbed{
		Title:       dailyDeed.Title,
		Description: embedDesc,
		Image: &discordgo.MessageEmbedImage{
			URL: dailyDeed.Image.Desktop,
		},
		Type: discordgo.EmbedTypeImage,
	}
	return &responseEmbed, nil
}
