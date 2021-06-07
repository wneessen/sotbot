package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
	"time"
)

// Get the daily deed
func GetDailyDeed(h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetDailyDeed",
	})
	dailyDeed, err := api.GetDailyDeed(h, u.RatCookie)
	if err != nil {
		l.Errorf("An error occurred fetching daily deeed data: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	embedDesc := fmt.Sprintf("%v\n\nStart date: %v\nEnd date: %v", dailyDeed.BodyText,
		dailyDeed.StartDate.Time().Format(time.RFC1123),
		dailyDeed.EndDate.Time().Format(time.RFC1123))
	embedFoot := &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("%v %v is waiting for you. Go for it!", dailyDeed.EntitlementRewardValue,
			dailyDeed.EntitlementRewardCurrency),
	}
	if dailyDeed.CompletedAt.Time().Unix() > 0 {
		embedFoot.Text = fmt.Sprintf("You already completed the deed on %v and earned %v %v",
			dailyDeed.CompletedAt.Time().Format(time.RFC1123), dailyDeed.EntitlementRewardValue,
			dailyDeed.EntitlementRewardCurrency)
	}
	responseEmbed := discordgo.MessageEmbed{
		Title:       dailyDeed.Title,
		Description: embedDesc,
		Image: &discordgo.MessageEmbedImage{
			URL: dailyDeed.Image.Desktop,
		},
		Footer: embedFoot,
		Type:   discordgo.EmbedTypeImage,
	}
	return &responseEmbed, nil
}
