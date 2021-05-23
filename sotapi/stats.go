package sotapi

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type UserStats struct {
	KrakenDefeated      ApiStringInt `json:"Combat_Kraken_Defeated"`
	MegalodonEncounters ApiStringInt `json:"Player_TinyShark_Spawned"`
	ChestsHandedIn      ApiStringInt `json:"Chests_HandedIn_Total"`
	ShipsSunk           ApiStringInt `json:"Combat_Ships_Sunk"`
	VomitedTotal        ApiStringInt `json:"Vomited_Total"`
}

type ApiOverview struct {
	Stats UserStats `json:"stats"`
}

func GetStats(hc *http.Client, rc string) (UserStats, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetStats",
	})
	apiUrl := "https://www.seaofthieves.com/api/profilev2/overview"

	l.Debugf("Fetching user season progress from API...")
	var userOverview ApiOverview
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, rc, "")
	if err != nil {
		return UserStats{}, err
	}
	if err := json.Unmarshal(httpResp, &userOverview); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return UserStats{}, err
	}

	return userOverview.Stats, nil
}
