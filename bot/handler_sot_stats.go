package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/sotapi"
)

// Just a test handler
func (b *Bot) SotStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.SotStats",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.Content == "!stats" {
		l.Debugf("Received '!stats' request from user %v", m.Author.Username)
		userObj, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			l.Errorf("Database lookup failed: %v", err)
			return
		}
		if userObj.ID <= 0 {
			replyMsg := "Sorry but your are not a registered user."
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		userRatCookie := database.UserGetPrefString(b.Db, userObj.ID, "rat_cookie")
		if userRatCookie == "" {
			replyMsg := "Sorry but you have no RAT cookie set. Try !setrat in the DMs"
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		userStats, err := sotapi.GetStats(b.HttpClient, userRatCookie)
		if err != nil {
			l.Errorf("An error occured fetching user stats: %v", err)
			replyMsg := fmt.Sprintf("Sorry but there was an error fetching your user stats"+
				" from the SoT API: %v", err)
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		statsResponse := fmt.Sprintf("During your journeys on the Sea of Thieves, so far, you defeated %d "+
			"kraken, had %d encounters with a Megalodon, handed in %d chests, sank %d other ships and vomited "+
			"%d times. Good job!", userStats.KrakenDefeated, userStats.MegalodonEncounters, userStats.ChestsHandedIn,
			userStats.ShipsSunk, userStats.VomitedTotal)
		AnswerUser(s, m, statsResponse, m.Author.Mention())
	}
}
