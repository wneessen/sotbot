package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) DevTestHandler(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	l := log.WithFields(log.Fields{
		"action": "handler.DevTestHandler",
	})
	if m.Activities != nil {
		for _, curActivity := range m.Activities {
			l.Debugf("%v started playing %v", m.User.ID, curActivity.Name)
		}
	}
}
