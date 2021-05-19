package bot

import (
	log "github.com/sirupsen/logrus"
)

func (b *Bot) Announce(m string) {
	l := log.WithFields(log.Fields{
		"action": "bot.Announce",
	})
	if b.AnnounceChan != nil {
		_, err := b.Session.ChannelMessageSend(b.AnnounceChan.ID, m)
		if err != nil {
			l.Errorf("Failed to make announcement: %v", err)
		}
	}
}
