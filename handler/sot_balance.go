package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
	"net/http"
)

// Get current SoT balance
func GetSotBalance(d *gorm.DB, h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotBalance",
	})

	var emFields []*discordgo.MessageEmbedField
	_ = u.UpdateSotBalance(d, h)
	userBalance, err := database.GetBalance(d, u.UserInfo.ID)
	if err != nil {
		l.Errorf("Database SoT balance lookup failed: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	p := message.NewPrinter(language.German)
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "\U0001F7E1 Gold ",
		Value:  fmt.Sprintf("📈 **%v** ", p.Sprintf("%d", userBalance.Gold)),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "🔵 Doubloons ",
		Value:  fmt.Sprintf("📈 **%v** ", p.Sprintf("%d", userBalance.Doubloons)),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "💰 Ancient Coins ",
		Value:  fmt.Sprintf("📈 **%v** ", p.Sprintf("%d", userBalance.AncientCoins)),
		Inline: true,
	})

	responseEmbed := &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  fmt.Sprintf("Current Sea of Thieves balance of user @%v", u.AuthorName),
		Fields: emFields,
	}
	return responseEmbed, nil
}
