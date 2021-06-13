package handler

import (
	"fmt"
	"github.com/wneessen/sotbot/cache"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
)

func GetTraderoutes(hc *http.Client, d *gorm.DB) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetTraderoutes",
	})

	var traderoutes api.Traderoutes
	var err error
	fromCache := true
	if err := cache.Read("traderoutes", &traderoutes, d); err != nil {
		l.Errorf("Failed to read traderoutes from DB cache: %v", err)
		fromCache = false
	}

	if fromCache {
		if traderoutes.ValidThru.Unix() < time.Now().Unix() {
			fromCache = false
		}
		l.Debugf("Cache traderoutes still valid. Using cached version")
	}

	if !fromCache {
		traderoutes, err = api.GetTraderoutes(hc)
		if err != nil {
			l.Errorf("An error occurred fetching traderoutes: %v", err)
			return &discordgo.MessageEmbed{}, err
		}
		if err := cache.Store("traderoutes", traderoutes, d); err != nil {
			l.Errorf("Failed to store traderoutes in DB cache: %v", err)
		}
	}

	var respondOutposts []*discordgo.MessageEmbedField
	for _, v := range traderoutes.Routes {
		respondOutposts = append(respondOutposts, &discordgo.MessageEmbedField{
			Name:   v.Outpost,
			Value:  fmt.Sprintf("⬆️ **%v** \n ⬇️ **%v**", strings.Title(v.Surplus), strings.Title(v.Sought)),
			Inline: true,
		})
	}

	responseEmbed := discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Trade Routes",
		Description: fmt.Sprintf("for %v", traderoutes.Dates),
		Footer:      &discordgo.MessageEmbedFooter{Text: "Source http://maps.seaofthieves.rarethief.com/"},
		Fields:      respondOutposts,
	}
	return &responseEmbed, nil
}
