package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/user"
	"runtime"
)

// Let's the bot tell you the current date/time when requested via !time command
func (b *Bot) TellMemUsage(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.TellMemUsage",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.Content == "!mem" {
		l.Debugf("Received '!mem' request from user %v", m.Author.Username)

		if !user.IsAdmin(s, m.Author.ID, m.ChannelID) {
			return
		}

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		memResponse := fmt.Sprintf(
			"Memory allocated: %v MiB\nTotal allocated: %v MiB\nSys Memory allocated: %v MiB\n"+
				"Number of GCs: %v",
			(memStats.Alloc / 1024 / 1024), (memStats.TotalAlloc / 1024 / 1024), (memStats.Sys / 1024 / 1024),
			memStats.NumGC)
		AnswerUser(s, m, "\n`"+memResponse+"`", "")
	}
}
