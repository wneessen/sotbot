package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/version"
)

func (b *Bot) BotReadyHandler(s *discordgo.Session, ev *discordgo.Ready) {
	l := log.WithFields(log.Fields{
		"action":    "handler.BotReadyHandler",
		"sessionID": ev.SessionID,
	})
	l.Debugf("Bot reached the 'ready' state!")

	usd := &discordgo.UpdateStatusData{Status: "online"}
	usd.Activities = make([]*discordgo.Activity, 1)
	usd.Activities[0] = &discordgo.Activity{
		Name: fmt.Sprintf("SoTBot v%v", version.Version),
		Type: discordgo.ActivityTypeGame,
		URL:  "https://github.com/wneessen/sotbot",
	}

	err := s.UpdateStatusComplex(*usd)
	if err != nil {
		l.Errorf("Failed to set status: %v", err)
	}
}
