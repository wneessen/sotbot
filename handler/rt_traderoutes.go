package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
)

func GetTraderoutes(hc *http.Client) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetTraderoutes",
	})
	traderoutes, err := api.GetTraderoutes(hc)
	if err != nil {
		l.Errorf("An error occurred fetching traderoutes: %v", err)
		return &discordgo.MessageEmbed{}, err
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
