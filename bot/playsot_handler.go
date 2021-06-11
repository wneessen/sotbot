package bot

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/tools"
	"github.com/wneessen/sotbot/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

// UserPlaysSot monitors the presence updates of the guild users. If a user
// starts playing SoT, it will fetch the current balance from the SoT API and
// store it in the database. Once the user finished their game, it will redo
// the same action and present a difference to the announce channel
func (b *Bot) UserPlaysSot(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	l := log.WithFields(log.Fields{
		"action": "bot.UserPlaySot",
	})

	userObj, err := user.NewUser(b.Db, b.Config, m.User.ID)
	if err != nil {
		l.Errorf("Failed to create new user object: %v", err)
		return
	}
	userObj.AuthorName = m.User.Username
	userObj.Mention = m.User.Mention()

	if !userObj.IsRegistered() {
		return
	}
	discordUser, err := s.User(m.User.ID)
	if err != nil {
		l.Errorf("Discord user lookup failed: %v", err)
		return
	}

	// User started an activity, SoT maybe?
	if len(m.Activities) > 0 {
		for _, curActivity := range m.Activities {
			if curActivity.Name == "Sea of Thieves" {
				l.Debugf("%v started playing SoT. Updating balance...", discordUser.Username)
				_ = userObj.UpdateSotBalance(b.Db, b.HttpClient)
				userStats, err := api.GetStats(b.HttpClient, userObj.RatCookie)
				if err == nil {
					userStats.ShipsSunk = userStats.ShipsSunk - 1
					var statsString bytes.Buffer
					gobEnc := gob.NewEncoder(&statsString)
					if err := gobEnc.Encode(userStats); err != nil {
						l.Errorf("Failed to serialize user stats.")
					}
					statsBase64 := base64.StdEncoding.EncodeToString(statsString.Bytes())
					if err := database.UserSetPref(b.Db, userObj.UserInfo.ID, "sot_stats",
						statsBase64); err != nil {
						l.Errorf("Failed to store user stats in DB: %v", err)
					}
				}

				if err := database.UserSetPref(b.Db, userObj.UserInfo.ID, "playing_sot",
					time.Now().Format(time.RFC3339)); err != nil {
					l.Errorf("Failed to update user status in database: %v", err)
				}
			}
		}
	}

	// User might have stopped an activity
	if len(m.Activities) == 0 {
		userWasPlaying := database.UserGetPrefString(b.Db, userObj.UserInfo.ID, "playing_sot")
		if userWasPlaying != "" {
			l.Debugf("%v stopped playing SoT. Updating balance...", discordUser.Username)
			userStartedPlaying, err := time.Parse(time.RFC3339, userWasPlaying)
			if err != nil {
				l.Errorf("Could not parse user start playing-time. Skipping announcement.")
				return
			}
			_ = userObj.UpdateSotBalance(b.Db, b.HttpClient)

			if err := database.UserDelPref(b.Db, userObj.UserInfo.ID, "playing_sot"); err != nil {
				l.Errorf("Failed to delete user status in database: %v", err)
			}

			userBalance, err := database.GetBalance(b.Db, userObj.UserInfo.ID)
			if err != nil {
				l.Errorf("Failed to read user balance from DB: %v", err)
				return
			}

			userStats, err := api.GetStats(b.HttpClient, userObj.RatCookie)
			if err != nil {
				l.Errorf("Failed to fetch user stats from SoT API: %v", err)
				return
			}

			statsChanged := false
			var oldStats api.UserStats
			oldStatsObjBase64 := database.UserGetPrefString(b.Db, userObj.UserInfo.ID, "sot_stats")
			if err := database.UserDelPref(b.Db, userObj.UserInfo.ID, "sot_stats"); err != nil {
				l.Errorf("Failed to delete user stats in database: %v", err)
			}
			if oldStatsObjBase64 != "" {
				var statsString bytes.Buffer
				oldStatsObjString, err := base64.StdEncoding.DecodeString(oldStatsObjBase64)
				if err != nil {
					l.Errorf("Failed to decode base64 string: %v", err)
				}
				statsString.Write(oldStatsObjString)
				gobDec := gob.NewDecoder(&statsString)
				if err := gobDec.Decode(&oldStats); err != nil {
					l.Errorf("Failed to serialize old user stats.")
				}
				statsChanged = tools.CompareSotStats(userStats, oldStats)
			}

			if userBalance.LastUpdated <= userStartedPlaying.Unix() && !statsChanged {
				l.Debugf("User balance seems to not have changed during game play")
				return
			}

			p := message.NewPrinter(language.German)
			if b.Config.GetBool("sot_play_announce") && b.AnnounceChan != nil {
				balDiff := database.GetBalanceDifference(b.Db, userObj.UserInfo.ID)
				balMsg := fmt.Sprintf("earned/spent: %v gold, %v doubloons and %v ancient coins.",
					p.Sprintf("%d", balDiff.Gold), p.Sprintf("%d", balDiff.Doubloons),
					p.Sprintf("%d", balDiff.AncientCoins))
				statMsg := fmt.Sprintf("defeated %v Kraken, encountered %v Megalodon(s), handed in %v chest(s), "+
					"sank %v ship(s) and vomited %v time(s)",
					p.Sprintf("%d", userStats.KrakenDefeated-oldStats.KrakenDefeated),
					p.Sprintf("%d", userStats.MegalodonEncounters-oldStats.MegalodonEncounters),
					p.Sprintf("%d", userStats.ChestsHandedIn-oldStats.ChestsHandedIn),
					p.Sprintf("%d", userStats.ShipsSunk-oldStats.ShipsSunk),
					p.Sprintf("%d", userStats.VomitedTotal-oldStats.VomitedTotal))

				var msg string
				if userBalance.LastUpdated >= userStartedPlaying.Unix() && statsChanged {
					if balDiff.Gold != 0 || balDiff.AncientCoins != 0 || balDiff.Doubloons != 0 {
						msg = fmt.Sprintf("On their last trip to the Sea of Thieves, %v %v"+
							"They also %v", userObj.Mention, balMsg, statMsg)
					} else {
						msg = fmt.Sprintf("On their last trip to the Sea of Thieves, %v didn't change "+
							"their balance, but %v", userObj.Mention, statMsg)
					}
				}
				if userBalance.LastUpdated >= userStartedPlaying.Unix() && !statsChanged {
					if balDiff.Gold != 0 || balDiff.AncientCoins != 0 || balDiff.Doubloons != 0 {
						msg = fmt.Sprintf("On their last trip to the Sea of Thieves, %v %v",
							userObj.Mention, balMsg)
					}
				}
				if userBalance.LastUpdated <= userStartedPlaying.Unix() && statsChanged {
					if balDiff.Gold != 0 || balDiff.AncientCoins != 0 || balDiff.Doubloons != 0 {
						msg = fmt.Sprintf("On their last trip to the Sea of Thieves, %v %v",
							userObj.Mention, statMsg)
					}
				}
				if msg != "" {
					response.Announce(s, b.AnnounceChan.ID, msg)
				}
			}
		}
	}
}
