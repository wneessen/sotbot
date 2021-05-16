package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) BotReadyHandler(s *discordgo.Session, ev *discordgo.Ready) {
	l := log.WithFields(log.Fields{
		"action":    "handler.BotReadyHandler",
		"sessionID": ev.SessionID,
	})
	l.Debugf("Bot reached the 'ready' state!")

	if err := s.UpdateListeningStatus("your commands"); err != nil {
		l.Errorf("Failed to set status: %v", err)
	}
}
