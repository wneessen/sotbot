package handler

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/audio"
)

func PlaySound(v []*discordgo.VoiceState, s *discordgo.Session, ab [][]byte, u, g string) error {
	l := log.WithFields(log.Fields{
		"action": "handler.PlaySound",
	})

	// Look for the message sender in that guild's current voice states.
	for _, vs := range v {
		if vs.UserID == u {
			vc, err := s.ChannelVoiceJoin(g, vs.ChannelID, false, true)
			if err != nil {
				l.Errorf("Failed to join voice chat: %v", err)
				_ = vc.Disconnect()
				return err
			}
			audio.PlayAudio(vc, ab)
			if err := vc.Disconnect(); err != nil {
				l.Errorf("Failed to disconnect voice chat: %v", err)
				return err
			}
		}
	}

	return nil
}
