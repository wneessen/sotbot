package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"time"
)

// PlaySound plays the current buffer to the provided channel.
func (b *Bot) PlayAudio(vc *discordgo.VoiceConnection, a string) error {
	l := log.WithFields(log.Fields{
		"action": "bot.PlayAudio",
	})
	if a == "" {
		return fmt.Errorf("Audio file %q not found in config", a)
	}
	time.Sleep(250 * time.Millisecond)
	if err := vc.Speaking(true); err != nil {
		l.Errorf("Failed to enable voice chat speaking mode: %v", err)
	}
	for _, buff := range *b.Audio[a].Buffer {
		vc.OpusSend <- buff
	}
	if err := vc.Speaking(false); err != nil {
		l.Errorf("Failed to disable voice chat speaking mode: %v", err)
	}
	time.Sleep(250 * time.Millisecond)

	return nil
}
