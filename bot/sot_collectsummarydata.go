package bot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/cache"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/user"
	"time"
)

func (b *Bot) CollectSummaryData() {
	l := log.WithFields(log.Fields{
		"action": "bot.CollectSummaryData",
	})

	userList, err := database.GetUsers(b.Db)
	if err != nil {
		l.Errorf("Failed to fetch user list from DB: %v", err)
		return
	}
	for _, curUser := range userList {
		userObj, err := user.NewUser(b.Db, b.Config, curUser.UserId)
		if err != nil {
			l.Errorf("Failed to create user object: %v", err)
			break
		}

		// Let's first check the last update time from the DB
		var lastCheck time.Time
		updateKey := fmt.Sprintf("summary_update_%v", userObj.UserInfo.UserId)
		if err := cache.Read(updateKey, &lastCheck, b.Db); err != nil {
			l.Errorf("Failed to read last summary update time for user %v from cache. Assuming first ever run.",
				userObj.UserInfo.UserId)
		}
		if time.Now().Unix()-lastCheck.Unix() < 86400 {
			l.Debugf("Last collection run for user %v was %v (less than 24h ago). Skipping for now.",
				userObj.UserInfo.UserId, lastCheck.String())
			return
		}

		if userObj.HasRatCookie() {
			userBalance, err := api.GetBalance(b.HttpClient, userObj.RatCookie)
			if err != nil {
				l.Errorf("Failed to fetch user balance from API: %v", err)
			} else {
				balKey := fmt.Sprintf("sot_balance_%v", userObj.UserInfo.UserId)
				if err := cache.Store(balKey, userBalance, b.Db); err != nil {
					l.Errorf("Failed to store user balance in cache: %v", err)
				}
			}
			userStats, err := api.GetStats(b.HttpClient, userObj.RatCookie)
			if err != nil {
				l.Errorf("Failed to fetch user balance from API: %v", err)
			} else {
				statsKey := fmt.Sprintf("sot_stats_%v", userObj.UserInfo.UserId)
				if err := cache.Store(statsKey, userStats, b.Db); err != nil {
					l.Errorf("Failed to store user stats in cache: %v", err)
				}
			}
			if err := cache.Store(updateKey, time.Now(), b.Db); err != nil {
				l.Errorf("Failed to store/update collection time in DB")
			}
			return
		}
		l.Errorf("User %v needs a summary update but seems to have no valid RAT cookie",
			userObj.UserInfo.UserId)
	}
}
