package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type UserAchievements struct {
	Sorted []SortedAchievement `json:"sorted"`
}

type SortedAchievement struct {
	Achievement Achievement `json:"achievement"`
}

type Achievement struct {
	Sort        int    `json:"Sort"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	MediaUrl    string `json:"MediaUrl"`
}

func GetLatestAchievement(hc *http.Client, rc string) (Achievement, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetLatestAchievement",
	})
	var userAchievements UserAchievements
	apiUrl := "https://www.seaofthieves.com/api/profilev2/achievements"

	l.Debugf("Fetching user achievements from API...")
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, rc, "")
	if err != nil {
		return Achievement{}, err
	}
	if err := json.Unmarshal(httpResp, &userAchievements); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return Achievement{}, err
	}

	return userAchievements.Sorted[0].Achievement, nil
}
