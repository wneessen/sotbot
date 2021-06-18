package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/cache"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
	"gorm.io/gorm"
	"net/http"
)

// Provide a daily summary
func GetSotSummary(h *http.Client, u *user.User, d *gorm.DB) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotSummary",
	})

	var oldBalance api.UserBalance
	keyName := fmt.Sprintf("sot_balance_%v", u.UserInfo.UserId)
	if err := cache.Read(keyName, &oldBalance, d); err != nil {
		return &discordgo.MessageEmbed{}, err
	}
	userBalance, err := api.GetBalance(h, u.RatCookie)
	if err != nil {
		l.Errorf("Failed to fetch user balance from API: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	summaryData := make(map[string]int64)
	summaryData["1_Gold"] = int64(userBalance.Gold) - int64(oldBalance.Gold)
	summaryData["2_Doubloon"] = int64(userBalance.Doubloons) - int64(oldBalance.Doubloons)
	summaryData["3_AncientCoin"] = int64(userBalance.AncientCoins) - int64(oldBalance.AncientCoins)

	var oldStats api.UserStats
	keyName = fmt.Sprintf("sot_stats_%v", u.UserInfo.UserId)
	if err := cache.Read(keyName, &oldStats, d); err != nil {
		return &discordgo.MessageEmbed{}, err
	}
	userStats, err := api.GetStats(h, u.RatCookie)
	if err != nil {
		l.Errorf("An error occurred fetching user stats: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	summaryData["4_Kraken"] = int64(userStats.KrakenDefeated) - int64(oldStats.KrakenDefeated)
	summaryData["5_Megalodon"] = int64(userStats.MegalodonEncounters) - int64(oldStats.MegalodonEncounters)
	summaryData["6_Chest"] = int64(userStats.ChestsHandedIn) - int64(oldStats.ChestsHandedIn)
	summaryData["7_Ship"] = int64(userStats.ShipsSunk) - int64(oldStats.ShipsSunk)
	summaryData["8_Vomit"] = int64(userStats.VomitedTotal) - int64(oldStats.VomitedTotal)

	emFields := response.FormatEmFields(summaryData)
	responseEmbed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("Daily Sea of Thieves summary for @%v", u.AuthorName),
		Description: "No changes happened since yesterday",
	}

	if len(emFields) > 0 {
		responseEmbed.Fields = emFields
		responseEmbed.Description = ""
	}

	return responseEmbed, nil
}
