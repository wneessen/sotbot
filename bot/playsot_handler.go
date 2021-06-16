package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/cache"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"sort"
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

	// Only track if we have an announce channel set
	if !b.Config.GetBool("sot_play_announce") || b.AnnounceChan == nil {
		l.Debugf("No SoT announce channel set or feature disabled")
		return
	}

	// User started an activity, SoT maybe?
	if len(m.Activities) > 0 {
		for _, curActivity := range m.Activities {
			if curActivity.Name == "Sea of Thieves" {
				l.Debugf("%v started playing SoT. Updating balance...", discordUser.Username)
				userBalance, err := api.GetBalance(b.HttpClient, userObj.RatCookie)
				if err == nil {
					balBase64, err := cache.SerializeObj(userBalance)
					if err != nil {
						l.Errorf("Failed to serialize user stats: %v", err)
					}
					if err == nil {
						if err := database.UserSetPref(b.Db, userObj.UserInfo.ID, "sot_balance",
							balBase64); err != nil {
							l.Errorf("Failed to store user stats in DB: %v", err)
						}
					}
				}

				userStats, err := api.GetStats(b.HttpClient, userObj.RatCookie)
				if err == nil {
					statsBase64, err := cache.SerializeObj(userStats)
					if err != nil {
						l.Errorf("Failed to serialize user stats: %v", err)
					}
					if err == nil {
						if err := database.UserSetPref(b.Db, userObj.UserInfo.ID, "sot_stats",
							statsBase64); err != nil {
							l.Errorf("Failed to store user stats in DB: %v", err)
						}
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
		if err := database.UserDelPref(b.Db, userObj.UserInfo.ID, "playing_sot"); err != nil {
			l.Errorf("Failed to delete user status in database: %v", err)
		}
		if userWasPlaying == "" {
			return
		}

		// Wait for some time to present the voyage statistics
		go func(t time.Time) {
			time.Sleep(time.Minute * 1)
			userIsPlaying := database.UserGetPrefString(b.Db, userObj.UserInfo.ID, "playing_sot")
			if userIsPlaying != "" {
				l.Debugf("%v apparently resumed playing...", discordUser.Username)
				return
			}
			l.Debugf("%v stopped playing SoT...", discordUser.Username)
			if err := database.UserDelPref(b.Db, userObj.UserInfo.ID, "playing_sot"); err != nil {
				l.Errorf("Failed to delete user status in database: %v", err)
			}

			voyageStats := make(map[string]int64, 0)
			userStartedPlaying, err := time.Parse(time.RFC3339, userWasPlaying)
			if err != nil {
				l.Errorf("Could not parse user start playing-time. Skipping announcement.")
				return
			}
			playTime := t.Unix() - userStartedPlaying.Unix()
			if playTime < 5 {
				l.Debugf("%v played less than 3 minutes (%v seconds). There is no chance of any change.",
					discordUser.Username, playTime)
				return
			}
			_ = userObj.UpdateSotBalance(b.Db, b.HttpClient)

			// Compare balance
			var oldBalance api.UserBalance
			userBalance, err := database.GetBalance(b.Db, userObj.UserInfo.ID)
			if err != nil {
				l.Errorf("Failed to fetch SoT user balance: %v", err)
			}
			oldBalanceObjBase64 := database.UserGetPrefString(b.Db, userObj.UserInfo.ID, "sot_balance")
			if err := database.UserDelPref(b.Db, userObj.UserInfo.ID, "sot_balance"); err != nil {
				l.Errorf("Failed to delete user balance in database: %v", err)
			}
			if oldBalanceObjBase64 != "" {
				if err := cache.DeserializeObj(oldBalanceObjBase64, &oldBalance); err != nil {
					l.Errorf("Failed to deserialize old user stats: %v", err)
					return
				}
			}
			voyageStats["1_Gold"] = int64(userBalance.Gold) - int64(oldBalance.Gold)
			voyageStats["2_Doubloon"] = int64(userBalance.Doubloons) - int64(oldBalance.Doubloons)
			voyageStats["3_AncientCoin"] = int64(userBalance.AncientCoins) - int64(oldBalance.AncientCoins)

			// Compare user stats
			var oldStats api.UserStats
			userStats, err := api.GetStats(b.HttpClient, userObj.RatCookie)
			if err != nil {
				l.Errorf("Failed to fetch user stats from SoT API: %v", err)
			}
			oldStatsObjBase64 := database.UserGetPrefString(b.Db, userObj.UserInfo.ID, "sot_stats")
			if err := database.UserDelPref(b.Db, userObj.UserInfo.ID, "sot_stats"); err != nil {
				l.Errorf("Failed to delete user stats in database: %v", err)
			}
			if oldBalanceObjBase64 != "" {
				if err := cache.DeserializeObj(oldStatsObjBase64, &oldStats); err != nil {
					l.Errorf("Failed to deserialize old user stats: %v", err)
					return
				}
			}
			voyageStats["4_Kraken"] = int64(userStats.KrakenDefeated) - int64(oldStats.KrakenDefeated)
			voyageStats["5_Megalodon"] = int64(userStats.MegalodonEncounters) - int64(oldStats.MegalodonEncounters)
			voyageStats["6_Chest"] = int64(userStats.ChestsHandedIn) - int64(oldStats.ChestsHandedIn)
			voyageStats["7_Ship"] = int64(userStats.ShipsSunk) - int64(oldStats.ShipsSunk)
			voyageStats["8_Vomit"] = int64(userStats.VomitedTotal) - int64(oldStats.VomitedTotal)

			// Prepare the output
			p := message.NewPrinter(language.German)
			var emFields []*discordgo.MessageEmbedField
			var keyNames []string
			for k := range voyageStats {
				keyNames = append(keyNames, k)
			}
			sort.Strings(keyNames)
			for _, k := range keyNames {
				v := voyageStats[k]
				if v != 0 {
					emFields = append(emFields, &discordgo.MessageEmbedField{
						Name: fmt.Sprintf("%v %v", response.Icon(k), response.IconKey(k)),
						Value: fmt.Sprintf("%v**%v** %v", response.BalanceIcon(k, v),
							p.Sprintf("%d", v), response.IconValue(k)),
						Inline: true,
					})
				}
			}

			if len(emFields) > 0 {
				// Response with the Embed
				responseEmbed := &discordgo.MessageEmbed{
					Type:   discordgo.EmbedTypeRich,
					Title:  fmt.Sprintf("Sea of Thieves voyage summary for @%v", discordUser.Username),
					Fields: emFields,
				}
				response.Embed(s, b.AnnounceChan.ID, responseEmbed)
			}
		}(time.Now())
	}
}
