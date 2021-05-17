package bot

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

func SetLogLevel(l string) {
	if l == "" {
		return
	}
	switch strings.ToLower(l) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
