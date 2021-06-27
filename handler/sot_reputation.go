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

	thumbUrls := map[string]string{
		"athena":   "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/c/c2/Athena%27s_Fortune_icon.png/revision/latest/scale-to-width-down/152?cb=20200512034100",
		"bilge":    "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/5/5e/Bilge_rat_adventures.png/revision/latest/scale-to-width-down/250?cb=20180710203959",
		"hoarder":  "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/f/f1/Gold_Hoarders_icon.png/revision/latest/scale-to-width-down/152?cb=20200512034953",
		"hunter":   "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/b/b1/The_Hunter%27s_Call_icon.png/revision/latest/scale-to-width-down/152?cb=20200920220917",
		"merchant": "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/9/9c/Merchant_Alliance_icon.png/revision/latest/scale-to-width-down/152?cb=20200512032051",
		"order":    "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/0/02/Order_of_Souls_icon.png/revision/latest/scale-to-width-down/152?cb=20200512154422",
		"reaper":   "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/6/61/Reaper%27s_Bones_icon.png/revision/latest/scale-to-width-down/152?cb=20200512033231",
		"seadog":   "https://static.wikia.nocookie.net/seaofthieves_gamepedia/images/3/3c/Sea_Dogs_icon.png/revision/latest/scale-to-width-down/152?cb=20200919191401",
	}

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
			URL: thumbUrls[factionName],
		},
		Fields: emFields,
	}
	return responseEmbed, nil
}
