package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
)

// Just a test handler
func GetSotSeasonProgress(h *http.Client, u *user.User) (string, bool, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotSeasonProgress",
	})

	userAchievement, err := api.GetSeasonProgress(h, u.RatCookie)
	if err != nil {
		l.Errorf("An error occured fetching user progress: %v", err)
		return "", false, err
	}
	responseMsg := fmt.Sprintf("You are currently sailing in %v. Your renown level is %v%% (Tier: %d). "+
		"Of the total amount of %d season challanges, so far, you completed %d.", userAchievement.SeasonTitle,
		fmt.Sprintf("%.1f", userAchievement.LevelProgress), userAchievement.Tier,
		userAchievement.TotalChallenges, userAchievement.CompletedChallenges)
	return responseMsg, true, nil
}
