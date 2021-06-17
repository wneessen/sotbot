package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"net/http"
	"strings"
)

func UrbanDict(h *http.Client, w string) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.UrbanDict",
	})

	urbanDict, err := api.GetUrbanDict(h, w)
	if err != nil {
		l.Errorf("Could not fetch UD response: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	urbanDict.Definition = strings.ReplaceAll(urbanDict.Definition, "[", "**")
	urbanDict.Definition = strings.ReplaceAll(urbanDict.Definition, "]", "**")
	urbanDict.Example = strings.ReplaceAll(urbanDict.Example, "[", "")
	urbanDict.Example = strings.ReplaceAll(urbanDict.Example, "]", "")

	responseEmbed := &discordgo.MessageEmbed{
		Title:       urbanDict.Word,
		Description: fmt.Sprintf("%v\n\nExample:\n`%v`", urbanDict.Definition[:1800], urbanDict.Example),
		Type:        discordgo.EmbedTypeArticle,
		URL:         urbanDict.PermaLink,
	}

	return responseEmbed, nil
}
