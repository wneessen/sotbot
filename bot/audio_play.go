package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// PlaySound plays the current buffer to the provided channel.
func (b *Bot) PlayAudio(s *discordgo.Session, g, c, a string) error {
	if a == "" {
		return fmt.Errorf("Audio file %q not found in config", a)
	}

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(g, c, false, true)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)
	if err := vc.Speaking(true); err != nil {
		return err
	}

	// Send the buffer data.
	for _, buff := range *b.Audio[a].Buffer {
		vc.OpusSend <- buff
	}

	if err := vc.Speaking(false); err != nil {
		return err
	}

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	if err := vc.Disconnect(); err != nil {
		return err
	}

	return nil
}
