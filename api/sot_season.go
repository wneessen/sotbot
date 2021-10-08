package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type Seasons []SeasonProgress

type SeasonProgress struct {
	LevelProgress       float64      `json:"LevelProgress"`
	Tier                int          `json:"Tier"`
	SeasonTitle         string       `json:"Title"`
	TotalChallenges     int          `json:"TotalChallenges"`
	CompletedChallenges int          `json:"CompleteChallenges"`
	Tiers               []SeasonTier `json:"Tiers"`
}

type SeasonTier struct {
	Number int    `json:"Number"`
	Title  string `json:"Title"`
}

func GetSeasonProgress(hc *http.Client, rc string) (SeasonProgress, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetSeasonProgress",
	})
	apiUrl := "https://www.seaofthieves.com/api/profilev2/seasons-progress"

	l.Debugf("Fetching user season progress from API...")
	var seasons Seasons
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, &rc, nil, false)
	if err != nil {
		return SeasonProgress{}, err
	}
	if err := json.Unmarshal(httpResp, &seasons); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return SeasonProgress{}, err
	}

	return seasons[len(seasons)-1], nil
}
