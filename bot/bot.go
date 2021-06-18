package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vascocosta/owm"
	"github.com/wneessen/sotbot/audio"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/httpclient"
	"github.com/wneessen/sotbot/random"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Bot struct {
	AuthToken    string
	Config       *viper.Viper
	Audio        map[string]Audio
	AudioMutex   *sync.Mutex
	HttpClient   *http.Client
	Db           *gorm.DB
	Session      *discordgo.Session
	AnnounceChan *discordgo.Channel
	TMDb         *tmdb.TMDb
	OwmClient    *owm.Client
	StartTime    time.Time
}

type Audio struct {
	Buffer *[][]byte
}

func NewBot(c *viper.Viper) Bot {
	l := log.WithFields(log.Fields{
		"action": "bot.NewBot",
	})
	authToken := c.GetString("authtoken")
	if authToken == "" {
		l.Errorf("AuthToken cannot be empty.")
		os.Exit(1)
	}
	bot := Bot{
		AuthToken:  authToken,
		Config:     c,
		AudioMutex: &sync.Mutex{},
		StartTime:  time.Now(),
	}
	bot.Audio = make(map[string]Audio)

	// Search and load audio files
	for af, fn := range c.GetStringMapString("audiofiles") {
		l.Debugf("Loading audio file %q as %q into memory", fn, af)

		audioBuffer := make([][]byte, 0)
		if err := audio.LoadAudio("./media/audio/"+fn, &audioBuffer); err != nil {
			l.Errorf("Failed to load audio file into memory: %v", err)
			break
		}
		bot.Audio[af] = Audio{
			Buffer: &audioBuffer,
		}
	}

	// Connect to database file
	dbObj, err := database.ConnectDB(c.GetString("dbfile"),
		c.GetString("loglevel"))
	if err != nil {
		l.Errorf("Failed to load database file: %v", err)
		os.Exit(1)
	}
	bot.Db = dbObj

	// Create a HTTP client object
	hc, err := httpclient.NewHttpClient()
	if err != nil {
		l.Errorf("Failed to create HTTP client: %v", err)
		os.Exit(1)
	}
	bot.HttpClient = hc

	// Create a TMDB object
	tmdbApiKey := c.GetString("tmdb_api_key")
	if tmdbApiKey != "" {
		bot.TMDb = tmdb.Init(tmdb.Config{
			APIKey:   tmdbApiKey,
			Proxies:  nil,
			UseProxy: false,
		})
	}

	// Create OWM object
	owmApiKey := c.GetString("owm_api_key")
	if owmApiKey != "" {
		bot.OwmClient = owm.NewClient(owmApiKey)
	}

	// Create/Fetch encyption key
	if err := bot.GetEncryptionKey(); err != nil {
		l.Errorf("Failed to create/read encryption key: %v", err)
		os.Exit(1)
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
	b.Session = discordObj

	// Do we have an announcement channel configured?
	configAnnounce := b.Config.GetString("announcechan")
	if configAnnounce != "" {
		chanObj, err := b.Session.Channel(configAnnounce)
		if err != nil {
			l.Errorf("Failed to look up discord channel: %v", err)
		} else {
			b.AnnounceChan = chanObj
		}
	}

	// Add handlers
	b.Session.AddHandler(b.BotReadyHandler)
	b.Session.AddHandler(b.CommandHandler)
	b.Session.AddHandler(b.SlashCmdHandler)
	b.Session.AddHandler(b.UserPlaysSot)

	// What events do we wanna see?
	b.Session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildPresences

	// Open the websocket and begin listening.
	err = b.Session.Open()
	if err != nil {
		l.Errorf("Error opening discord session: %v", err)
		return
	}

	// Get list of registered slash commands
	registeredCmds, err := b.Session.ApplicationCommands(b.Session.State.User.ID, GuildID)
	if err != nil {
		l.Errorf("Failed to read slash command list from Discord server: %v", err)
	}

	// Register the slash commands (if not already registered)
	for _, slashCmd := range b.SlashCmdList() {
		isNew := true
		hasChanged := false
		for _, regCmd := range registeredCmds {
			if slashCmd.Name == regCmd.Name && slashCmd.Description == regCmd.Description {
				l.Debugf("Slash command %q already registered. Skipping.", regCmd.Name)
				isNew = false
				hasChanged = false
				break
			}
			if slashCmd.Name == regCmd.Name && slashCmd.Description != regCmd.Description {
				l.Debugf("Slash command %q changed. Editing...", regCmd.Name)
				isNew = false
				hasChanged = true
				slashCmd.ID = regCmd.ID
				break
			}
		}
		if isNew {
			go func(s *discordgo.ApplicationCommand) {
				randNum, _ := random.Number(2000)
				randNum += 1000
				randDelay, _ := time.ParseDuration(fmt.Sprintf("%dms", randNum))
				time.Sleep(randDelay)
				l.Debugf("[%v] Registering slash command...", s.Name)
				_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, GuildID, s)
				if err != nil {
					l.Errorf("[%v] Registration failed: %v", s.Name, err)
					return
				}
				l.Debugf("[%v] Registration completed", s.Name)
			}(slashCmd)
		}
		if hasChanged {
			go func(s *discordgo.ApplicationCommand) {
				randNum, _ := random.Number(2000)
				randNum += 1000
				randDelay, _ := time.ParseDuration(fmt.Sprintf("%dms", randNum))
				time.Sleep(randDelay)
				l.Debugf("[%v] Updating slash command...", s.Name)
				_, err = b.Session.ApplicationCommandEdit(b.Session.State.User.ID, GuildID, s.ID, s)
				if err != nil {
					l.Errorf("[%v] Update failed: %v", s.Name, err)
					return
				}
				l.Debugf("[%v] Update completed", s.Name)
			}(slashCmd)
		}
	}

	// We need a signal channel
	sc := make(chan os.Signal, 1)
	signal.Notify(sc)

	// We want timed events as well
	checkAuthTimer := time.NewTicker(time.Hour)
	defer checkAuthTimer.Stop()
	summaryTimer := time.NewTicker(time.Minute * 30)
	defer summaryTimer.Stop()

	// Wait here until CTRL-C or other term signal is received.
	l.Infof("Bot is ready and connected. Press CTRL-C to exit.")
	for {
		select {
		case rs := <-sc:
			if rs == syscall.SIGKILL ||
				rs == syscall.SIGABRT ||
				rs == syscall.SIGINT ||
				rs == syscall.SIGTERM {
				l.Warnf("received %q signal. Exiting.", rs)

				// Cleanly close down the Discord session.
				if err := b.Session.Close(); err != nil {
					l.Errorf("Failed to gracefully close discord session: %v", err)
				}

				os.Exit(0)
			}
		case <-checkAuthTimer.C:
			go b.CheckSotAuth()

		case <-summaryTimer.C:
			go b.CollectSummaryData()
		}
	}
}

func (b *Bot) ResetSlashCmds() {
	l := log.WithFields(log.Fields{
		"action": "bot.ResetSlashCmds",
	})
	l.Infof("Initializing bot for 'reset all slash commands' run...")

	discordObj, err := discordgo.New("Bot " + b.AuthToken)
	if err != nil {
		l.Errorf("Error creating discord session: %v", err)
		return
	}
	b.Session = discordObj

	// What events do we wanna see?
	b.Session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildPresences

	// Open the websocket and begin listening.
	err = b.Session.Open()
	if err != nil {
		l.Errorf("Error opening discord session: %v", err)
		return
	}

	// Get list of registered slash commands
	registeredCmds, err := b.Session.ApplicationCommands(b.Session.State.User.ID, GuildID)
	if err != nil {
		l.Errorf("Failed to read slash command list from Discord server: %v", err)
	}

	// Delete all registered slash commands
	for _, regCmd := range registeredCmds {
		l.Debugf("Removing slash command %q...", regCmd.Name)
		err = b.Session.ApplicationCommandDelete(b.Session.State.User.ID, GuildID, regCmd.ID)
		if err != nil {
			l.Errorf("Failed to delete slash command: %v", err)
			continue
		}
		l.Debugf("Slash command %q successfully removed...", regCmd.Name)
	}

	if err := b.Session.Close(); err != nil {
		l.Errorf("Failed to close discord connection: %v", err)
	}
	os.Exit(0)
}
