package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
	"net/http"
)

// Just a test handler
func GetSotStats(h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotStats",
	})
	userStats, err := api.GetStats(h, u.RatCookie)
	if err != nil {
		l.Errorf("An error occurred fetching user stats: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	statsData := make(map[string]int64)
	statsData["1_Kraken"] = int64(userStats.KrakenDefeated)
	statsData["2_Megalodon"] = int64(userStats.MegalodonEncounters)
	statsData["3_Chest"] = int64(userStats.ChestsHandedIn)
	statsData["4_Ship"] = int64(userStats.ShipsSunk)
	statsData["5_Vomit"] = int64(userStats.VomitedTotal)

	emFields := response.FormatEmFields(statsData)
	responseEmbed := &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  fmt.Sprintf("Current Sea of Thieves stats of user @%v", u.AuthorName),
		Fields: emFields,
	}
	return responseEmbed, nil
}
