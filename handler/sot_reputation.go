package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
	"regexp"
	"strings"
)

// Just a test handler
func GetSotReputation(h *http.Client, u *user.User, f string) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotReputation",
	})

	wrongFormatMsg := fmt.Sprintf("Incorrect request format. Usage: !reputation " +
		"<athena|bilge|hoarder|hunter|merchant|order|reaper|seadog>")
	validFaction, err := regexp.Compile(`^(athena|bilge|hoarder|hunter|merchant|order|reaper|seadog)$`)
	if err != nil {
		return &discordgo.MessageEmbed{}, err
	}
	if !validFaction.MatchString(f) {
		return &discordgo.MessageEmbed{}, fmt.Errorf(wrongFormatMsg)
	}
	validFactionMatch := validFaction.FindStringSubmatch(f)
	if len(validFactionMatch) < 1 {
		return &discordgo.MessageEmbed{}, fmt.Errorf(wrongFormatMsg)
	}
	factionName := strings.ToLower(validFactionMatch[0])
	userReputation, err := api.GetFactionReputation(h, u.RatCookie, factionName)
	if err != nil {
		l.Errorf("An error occurred fetching user progress: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	var emFields []*discordgo.MessageEmbedField
	if userReputation.Rank != "" {
		emFields = append(emFields, &discordgo.MessageEmbedField{
			Name:   "Rank",
			Value:  userReputation.Rank,
			Inline: false,
		})
	}
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Level",
		Value:  fmt.Sprintf("🌡️ %d", userReputation.Level),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "XP in the current level",
		Value:  fmt.Sprintf("☑️ %d/%d", userReputation.Xp, userReputation.NextLevel.XpRequired),
		Inline: true,
	})

	responseEmbed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("%v reputation summary for @%v", userReputation.Name, u.AuthorName),
		Description: fmt.Sprintf("*%q*", userReputation.Motto),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://github.com/wneessen/sotbot/raw/main/assets/reputation/%s.png", factionName),
		},
		Fields: emFields,
	}
	return responseEmbed, nil
}
