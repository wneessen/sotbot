package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/sotbot"
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
	// Read CLI flags
	flag.Usage = printHelp
	//confFile := flag.String("c", "foo", "Path to config file")
	botToken := flag.String("t", "", "Bot authentication token")
	flag.Parse()

	if *botToken == "" {
		fmt.Printf("Error: No bot authentication token given.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	bot := sotbot.NewBot(*botToken)
	bot.Run()
}
