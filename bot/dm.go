package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func DmUser(s *discordgo.Session, u string, msg string) {
	l := log.WithFields(log.Fields{
		"action": "bot.DmUser",
	})

	st, err := s.UserChannelCreate(u)
	if err != nil {
		l.Errorf("Failed to initiate DM channel with user: %v", err)
		return
	}

	_, err = s.ChannelMessageSend(st.ID, msg)
	if err != nil {
		l.Errorf("Failed to notify user: %v", err)
		return
	}
}
