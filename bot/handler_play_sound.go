package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (b *Bot) PlaySound(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.PlaySound",
	})
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!airhorn" {
		AnswerUser(s, m, "!airhorn is now !play airhorn (or any other registered sound)", m.Author.Mention())
		return
	}

	if strings.HasPrefix(m.Content, "!play") {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			l.Errorf("Could not find channel: %v", m.ChannelID)
			return
		}

		msgArray := strings.SplitN(m.Content, " ", 2)
		if len(msgArray) != 2 {
			AnswerUser(s, m, "Wrong format. Usage: !play <sound_name>", m.Author.Mention())
			return
		}

		if b.Audio[msgArray[1]].Buffer != nil {
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
					err = b.PlayAudio(vc, msgArray[1])
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

		AnswerUser(s, m, fmt.Sprintf("Sorry, I don't have a registered sound file %q", msgArray[1]),
			m.Author.Mention())
		return
	}
}
