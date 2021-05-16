package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wneessen/sotbot/bot"
	"os"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: true,
		DisableColors:          false,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05 -0700",
	})
}

func main() {
	flag.Usage = printHelp
	confDir := flag.String("c", "", "Add custom config path for bot.json file")
	flag.Parse()

	botConf := viper.New()
	botConf.SetConfigName("sotbot")
	botConf.SetConfigType("json")

	if *confDir != "" {
		botConf.AddConfigPath(*confDir)
	}
	botConf.AddConfigPath("$HOME/.sotbot")
	botConf.AddConfigPath("./config")
	if err := botConf.ReadInConfig(); err != nil {
		log.Errorf("Failed to read config file: %v", err)
		os.Exit(1)
	}

	botObj := bot.NewBot(botConf)
	botObj.Run()
}
