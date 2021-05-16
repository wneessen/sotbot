package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Bot struct {
	AuthToken  string
	Config     *viper.Viper
	Audio      map[string]Audio
	AudioMutex *sync.Mutex
}

type Audio struct {
	Buffer *[][]byte
}

func NewBot(c *viper.Viper) Bot {
	l := log.WithFields(log.Fields{
		"action": "bot.NewBot",
	})
	authToken := c.GetString("AuthToken")
	if authToken == "" {
		l.Errorf("AuthToken cannot be empty.")
		os.Exit(1)
	}
	bot := Bot{
		AuthToken:  authToken,
		Config:     c,
		AudioMutex: &sync.Mutex{},
	}
	bot.Audio = make(map[string]Audio)

	for af, fn := range c.GetStringMapString("AudioFiles") {
		l.Debugf("Loading audio file %q as %q into memory", fn, af)

		audioBuffer := make([][]byte, 0)
		if err := LoadAudio("./media/audio/"+fn, &audioBuffer); err != nil {
			l.Errorf("Failed to load audio file into memory: %v", err)
			break
		}
		bot.Audio[af] = Audio{
			Buffer: &audioBuffer,
		}
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
	discordObj.AddHandler(b.BotReadyHandler)
	discordObj.AddHandler(b.TellTime)
	discordObj.AddHandler(b.Airhorn)

	// What events do we wanna see?
	discordObj.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsDirectMessages

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
