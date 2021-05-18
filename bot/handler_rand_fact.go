package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/factapi"
)

func (b *Bot) RandomFact(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.RandomFact",
	})
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!fact" {
		randFact, err := factapi.GetRandFact(b.HttpClient)
		if err != nil {
			l.Errorf("Could not fetch random fact: %v", err)
		}
		AnswerUser(s, m, randFact.Text)
	}
}
