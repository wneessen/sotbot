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
		Fields:      emFields,
	}
	return responseEmbed, nil
}
