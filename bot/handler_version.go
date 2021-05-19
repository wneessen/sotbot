package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/version"
)

// Let the bot tell us it's version information
func (b *Bot) TellVersion(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.TellVersion",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.Content == "!version" {
		l.Debugf("Received '!version' request from user %v", m.Author.Username)
		returnMsg := fmt.Sprintf("I am SoTBot Version v%v (built on: %v, built by: %v)", version.Version,
			version.BuildDate, version.BuildUser)
		AnswerUser(s, m, returnMsg, m.Author.Mention())
	}
}
