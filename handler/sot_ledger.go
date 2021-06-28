package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/http"
	"regexp"
	"strings"
)

// GetSotLedger is a SoTBot handler that replies the requesting user with their current
// SoT ledger position in a specific faction/company
func GetSotLedger(h *http.Client, u *user.User, f string) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotLedger",
	})

	wrongFormatMsg := fmt.Sprintf("Incorrect request format. Usage: !ledger " +
		"<athena|hoarder|merchant|order|reaper>")
	var validFaction, err = regexp.Compile(`^(athena|hoarder|merchant|order|reaper)$`)
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
	userLedger, err := api.GetFactionLedger(h, u.RatCookie, factionName)
	if err != nil {
		l.Errorf("An error occurred fetching user progress: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	var emFields []*discordgo.MessageEmbedField
	p := message.NewPrinter(language.German)
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Faction/Company",
		Value:  userLedger.Name,
		Inline: false,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Title",
		Value:  userLedger.BandTitle,
		Inline: false,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Emissary value",
		Value:  p.Sprintf("💰 %d", userLedger.Score),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Ledger position",
		Value:  p.Sprintf("🌡️ %d", userLedger.Rank),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Next level in",
		Value:  p.Sprintf("📈 %d points", userLedger.ToNextRank),
		Inline: true,
	})

	responseEmbed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: fmt.Sprintf("Global ledger position for @%v", u.AuthorName),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://github.com/wneessen/sotbot/raw/main/assets/ledger/%s%d.png",
				factionName, 4-userLedger.Band),
		},
		Fields: emFields,
	}
	return responseEmbed, nil
}
