package sotapi

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type SeasonProgress struct {
	LevelProgress       float64 `json:"LevelProgress"`
	Tier                int     `json:"Tier"`
	SeasonTitle         string  `json:"Title"`
	TotalChallenges     int     `json:"TotalChallenges"`
	CompletedChallenges int     `json:"CompleteChallenges"`
}

func GetSeasonProgress(hc *http.Client, rc string) (SeasonProgress, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetSeasonProgress",
	})
	apiUrl := "https://www.seaofthieves.com/api/profilev2/seasons-progress"

	l.Debugf("Fetching user season progress from API...")
	var userProgress SeasonProgress
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, rc, "")
	if err != nil {
		return userProgress, err
	}
	if err := json.Unmarshal(httpResp, &userProgress); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return userProgress, err
	}

	return userProgress, nil
}
