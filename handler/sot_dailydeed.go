package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/cache"
	"github.com/wneessen/sotbot/user"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// Get the daily deed
func GetDailyDeed(h *http.Client, u *user.User, d *gorm.DB) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetDailyDeed",
	})

	var dailyDeed api.BountyList
	var err error
	fromCache := true
	if err := cache.Read("dailydeed", &dailyDeed, d); err != nil {
		l.Errorf("Failed to read daily deed from DB cache: %v", err)
		fromCache = false
	}

	if fromCache {
		if dailyDeed.EndDate.Unix() < time.Now().Unix() {
			fromCache = false
		}
		l.Debugf("Daily deed still valid. Using cached version")
	}

	if !fromCache {
		dailyDeed, err = api.GetDailyDeed(h, u.RatCookie)
		if err != nil {
			l.Errorf("An error occurred fetching daily deed: %v", err)
			return &discordgo.MessageEmbed{}, err
		}
		if err := cache.Store("dailydeed", dailyDeed, d); err != nil {
			l.Errorf("Failed to store traderoutes in DB cache: %v", err)
		}
	}

	embedDesc := fmt.Sprintf("%v\n\nStart date: %v\nEnd date: %v", dailyDeed.BodyText,
		dailyDeed.StartDate.Format(time.RFC1123),
		dailyDeed.EndDate.Format(time.RFC1123))
	embedFoot := &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("%v %v is waiting for you. Go for it!", dailyDeed.EntitlementRewardValue,
			dailyDeed.EntitlementRewardCurrency),
	}
	if dailyDeed.CompletedAt.Unix() > 0 {
		embedFoot.Text = fmt.Sprintf("You already completed the deed on %v and earned %v %v",
			dailyDeed.CompletedAt.Format(time.RFC1123), dailyDeed.EntitlementRewardValue,
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
