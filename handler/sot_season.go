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
func GetSotSeasonProgress(h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotSeasonProgress",
	})

	romanNumeral := map[int]string{
		1:  "I",
		2:  "II",
		3:  "III",
		4:  "IV",
		5:  "V",
		6:  "VI",
		7:  "VII",
		8:  "VIII",
		9:  "IX",
		10: "X",
	}

	userAchievement, err := api.GetSeasonProgress(h, u.RatCookie)
	if err != nil {
		l.Errorf("An error occurred fetching user progress: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	var emFields []*discordgo.MessageEmbedField
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Title",
		Value:  userAchievement.Tiers[userAchievement.Tier-1].Title,
		Inline: false,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Renown Level",
		Value:  fmt.Sprintf("🌡️ %.1f%%", userAchievement.LevelProgress),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Renown Tier",
		Value:  fmt.Sprintf("📜 %d", userAchievement.Tier),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Challenges",
		Value:  fmt.Sprintf("☑️ %d/%d completed", userAchievement.CompletedChallenges, userAchievement.TotalChallenges),
		Inline: true,
	})

	responseEmbed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("Season summary for @%v", u.AuthorName),
		Description: userAchievement.SeasonTitle,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://github.com/wneessen/sotbot/raw/main/assets/numerals/%s.png",
				romanNumeral[userAchievement.Tier]),
		},
		Fields: emFields,
	}
	return responseEmbed, nil
}
