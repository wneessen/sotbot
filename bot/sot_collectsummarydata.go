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

	var lastCheck time.Time
	if err := cache.Read("summary_update", &lastCheck, b.Db); err != nil {
		l.Errorf("Failed to read last summary update time from cache. Assuming first ever run.")
	}
	if time.Now().Unix()-lastCheck.Unix() < 86400 {
		l.Debugf("Last collection run was %v (less than 24h ago). Skipping for now.", lastCheck.String())
		return
	}

	l.Debugf("Collecting data...")
	if err := cache.Store("summary_update", time.Now(), b.Db); err != nil {
		l.Errorf("Failed to store/update collection time in DB")
	}

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
		}

	}
}
