package audio

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"time"
)

// PlaySound plays the current buffer to the provided channel.
func PlayAudio(vc *discordgo.VoiceConnection, b [][]byte) {
	l := log.WithFields(log.Fields{
		"action": "audio.PlayAudio",
	})
	time.Sleep(100 * time.Millisecond)
	if err := vc.Speaking(true); err != nil {
		l.Errorf("Failed to enable voice chat speaking mode: %v", err)
	}
	for _, buff := range b {
		vc.OpusSend <- buff
	}
	if err := vc.Speaking(false); err != nil {
		l.Errorf("Failed to disable voice chat speaking mode: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
}
