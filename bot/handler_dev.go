package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (b *Bot) DevTestHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.DevTestHandler",
	})

	if m.Content == "!test" {
		l.Debugf("Test invoked")

		reqUser, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			return
		}

		balDiff := database.GetBalanceDifference(b.Db, reqUser.ID)
		p := message.NewPrinter(language.German)
		msg := fmt.Sprintf("Since their last trip to the Sea of Thieves, %v earned/lost: %v gold, "+
			"%v doubloons and %v ancient coins.", m.Author.Mention(), p.Sprintf("%d", balDiff.Gold),
			p.Sprintf("%d", balDiff.Doubloons), p.Sprintf("%d", balDiff.AncientCoins))
		b.Announce(msg)
	}
}
