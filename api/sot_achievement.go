package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

// UserAchievements defines the structure of the SoT API response for the
// users latest achievements
type UserAchievements struct {
	Sorted []SortedAchievement `json:"sorted"`
}

// SortedAchievement is a subpart of the UserAchievements API response
type SortedAchievement struct {
	Achievement Achievement `json:"achievement"`
}

// Achievement is a subpart of the SortedAchievement API response
type Achievement struct {
	Sort        int    `json:"Sort"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	MediaUrl    string `json:"MediaUrl"`
}

// GetLatestAchievement calls the SoT achievements API endpoint and returns
// a Achievement struct
func GetLatestAchievement(hc *http.Client, rc string) (Achievement, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetLatestAchievement",
	})
	var userAchievements UserAchievements
	apiUrl := "https://www.seaofthieves.com/api/profilev2/achievements"

	l.Debugf("Fetching user achievements from API...")
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, &rc, nil, false)
	if err != nil {
		return Achievement{}, err
	}
	if err := json.Unmarshal(httpResp, &userAchievements); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return Achievement{}, err
	}

	return userAchievements.Sorted[0].Achievement, nil
}
