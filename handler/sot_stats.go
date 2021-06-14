package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
)

// Just a test handler
func GetSotStats(h *http.Client, u *user.User) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotStats",
	})
	userStats, err := api.GetStats(h, u.RatCookie)
	if err != nil {
		l.Errorf("An error occurred fetching user stats: %v", err)
		return "", err
	}
	responseMsg := fmt.Sprintf("During your journeys on the Sea of Thieves, so far, you defeated %d "+
		"kraken, had %d encounters with a Megalodon, handed in %d chests, sank %d other ships and vomited "+
		"%d times. Good times!", userStats.KrakenDefeated, userStats.MegalodonEncounters, userStats.ChestsHandedIn,
		userStats.ShipsSunk, userStats.VomitedTotal)
	return responseMsg, nil
}
