package sotbot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/handler"
	"os"
	"os/signal"
	"syscall"
)

type Bot struct {
	AuthToken string
}

func NewBot(t string) Bot {
	bot := Bot{
		AuthToken: t,
	}

	return bot
}

func (b *Bot) Run() {
	l := log.WithFields(log.Fields{
		"action": "bot.Run",
	})
	l.Infof("Initializing bot...")

	discordObj, err := discordgo.New("Bot " + b.AuthToken)
	if err != nil {
		l.Errorf("Error creating discord session: %v", err)
		return
	}

	// Add handlers
	discordObj.AddHandler(handler.BotReadyHandler)

	// What events do we wanna see?
	discordObj.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates

	// Open the websocket and begin listening.
	err = discordObj.Open()
	if err != nil {
		l.Errorf("Error opening discoard session: %v", err)
		return
	}

	// We need a signal channel
	sc := make(chan os.Signal, 1)
	signal.Notify(sc)

	// Wait here until CTRL-C or other term signal is received.
	l.Infof("Bot is ready and connected. Press CTRL-C to exit.")
	for {
		select {
		case rs := <-sc:
			if rs == syscall.SIGKILL ||
				rs == syscall.SIGABRT ||
				rs == syscall.SIGINT ||
				rs == syscall.SIGTERM {
				l.Infof("received %q signal. Exiting.", rs)

				// Cleanly close down the Discord session.
				if err := discordObj.Close(); err != nil {
					l.Errorf("Failed to gracefully close discord session: %v", err)
				}

				os.Exit(0)
			}
		}
	}
}
