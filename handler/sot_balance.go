package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/http"
	"sort"
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

	p := message.NewPrinter(language.German)
	var emFields []*discordgo.MessageEmbedField
	var keyNames []string
	for k := range balanceData {
		keyNames = append(keyNames, k)
	}
	sort.Strings(keyNames)
	for _, k := range keyNames {
		v := balanceData[k]
		if v != 0 {
			emFields = append(emFields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%v %v", response.Icon(k), response.IconKey(k)),
				Value:  fmt.Sprintf("**%v** %v", p.Sprintf("%d", v), response.IconValue(k)),
				Inline: true,
			})
		}
	}
	for len(emFields)%3 != 0 {
		emFields = append(emFields, &discordgo.MessageEmbedField{
			Value:  "\U0000FEFF",
			Name:   "\U0000FEFF",
			Inline: true,
		})
	}
	responseEmbed := &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  fmt.Sprintf("Current Sea of Thieves balance of user @%v", u.AuthorName),
		Fields: emFields,
	}
	return responseEmbed, nil
}
