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

// Get current SoT balance
func GetSotBalance(h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotBalance",
	})

	userBalance, err := api.GetBalance(h, u.RatCookie)
	if err != nil {
		l.Errorf("Failed to fetch user balance from API: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	balanceData := make(map[string]int64)
	balanceData["1_Gold"] = int64(userBalance.Gold)
	balanceData["2_Doubloon"] = int64(userBalance.Doubloons)
	balanceData["3_AncientCoin"] = int64(userBalance.AncientCoins)

	emFields := response.FormatEmFields(balanceData)
	responseEmbed := &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  fmt.Sprintf("Current Sea of Thieves balance of user @%v", u.AuthorName),
		Fields: emFields,
	}
	return responseEmbed, nil
}
