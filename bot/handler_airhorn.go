package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (b *Bot) Airhorn(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.Airhorn",
	})
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!airhorn") {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			l.Errorf("Could not find channel: %v", m.ChannelID)
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			l.Errorf("Could not find guild of channel: %v", c.GuildID)
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					l.Errorf("Failed to join voice chat: %v", err)
					_ = vc.Disconnect()
					return
				}
				b.AudioMutex.Lock()
				err = b.PlayAudio(vc, "airhorn")
				if err != nil {
					l.Errorf("Failed to play audio: %v", err)
				}
				b.AudioMutex.Unlock()
				if err := vc.Disconnect(); err != nil {
					l.Errorf("Failed to disconnect voice chat: %v", err)
					return
				}

				return
			}
		}
	}
}
