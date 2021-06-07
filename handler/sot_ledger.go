package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/http"
	"regexp"
	"strings"
)

// Just a test handler
func GetSotLedger(h *http.Client, u *user.User, f string) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotLedger",
	})

	wrongFormatMsg := fmt.Sprintf("Incorrect request format. Usage: !ledger " +
		"<athena|hoarder|merchant|order|reaper>")
	var validFaction, err = regexp.Compile(`^(athena|hoarder|merchant|order|reaper)$`)
	if err != nil {
		return "", err
	}
	if !validFaction.MatchString(f) {
		return wrongFormatMsg, nil
	}
	validFactionMatch := validFaction.FindStringSubmatch(f)
	if len(validFactionMatch) < 1 {
		return wrongFormatMsg, nil
	}
	userLedger, err := api.GetFactionLedger(h, u.RatCookie, strings.ToLower(validFactionMatch[0]))
	if err != nil {
		l.Errorf("An error occurred fetching user progress: %v", err)
		return "", err
	}

	p := message.NewPrinter(language.German)
	responseMsg := fmt.Sprintf("You current global ledger rank within the **%v** faction is: **%v**. Your current"+
		" emissary value is **%v**, which results in position %v on the leaderboard. To reach the next ledger rank, "+
		"you'll need to increase your faction's emissary value by **%v points**.",
		userLedger.Name, userLedger.BandTitle, p.Sprintf("%d", userLedger.Score),
		p.Sprintf("%d", userLedger.Rank), p.Sprintf("%d", userLedger.ToNextRank))

	return responseMsg, nil
}
