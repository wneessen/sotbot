package bot

import (
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/random"
)

func RandomArrr() string {
	l := log.WithFields(log.Fields{
		"action": "bot.RandomArrr",
	})
	var pirateWords = []string{
		"Arrr", "Yarrr", "Arrgh", "Ahoy", "Garrr", "Yo-ho-ho",
	}
	wordCount := len(pirateWords)
	randNum, err := random.Number(wordCount)
	if err != nil {
		l.Errorf("Failed to generate random number: %v", err)
		return ""
	}

	return pirateWords[randNum]
}
