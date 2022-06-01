package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type Data struct {
	Adventures []Adventure `json:"Arcs"`
}

type Adventure struct {
	Meta              AdventureMeta         `json:"Meta"`
	GroupGoalProgress *AdvGroupGoalProgress `json:"GroupGoalProgress,omitempty"`
}

type AdventureMeta struct {
	Title       string `json:"Title"`
	AdventureID string `json:"AdventureId"`
	WebCode     string `json:"WebCode"`
}

type AdvGroupGoalProgress struct {
	LeadingGroupGoalID string    `json:"LeadingGroupGoalId"`
	Goals              []AdvGoal `json:"Goals"`
}

type AdvGoal struct {
	GroupGoalID string `json:"GroupGoalId"`
	State       string `json:"State"`
}

// GetAdventures calls the SoT adventures API endpoint and returns a Adventures struct
func GetAdventures(hc *http.Client, rc string) ([]Adventure, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetAdventures",
	})
	var data []Data
	apiUrl := "https://www.seaofthieves.com/api/profilev2/adventures"

	l.Debugf("Fetching adventures from API...")
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, &rc, nil, false)
	if err != nil {
		return []Adventure{}, err
	}
	if err := json.Unmarshal(httpResp, &data); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return []Adventure{}, err
	}

	return data[0].Adventures, nil
}
