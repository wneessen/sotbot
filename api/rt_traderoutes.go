package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
	"regexp"
)

// Traderoutes defines the structure for the rarethief.com trading route
// API response
type Traderoutes struct {
	Dates  string           `json:"trade_route_dates"`
	Routes map[string]Route `json:"routes"`
}

// Route defines the structure of the sub part of the API response that is
// a actual route
type Route struct {
	Outpost string `json:"outpost"`
	Sought  string `json:"sought_after"`
	Surplus string `json:"surplus"`
}

// GetTraderoutes fetches the currently active trading routes from the
// rarethief.com API and returns Traderoutes struct
func GetTraderoutes(hc *http.Client) (Traderoutes, error) {
	l := log.WithFields(log.Fields{
		"action": "rarethief.Traderoutes",
	})
	apiUrl := "https://maps.seaofthieves.rarethief.com/js/trade_routes.js"

	l.Debugf("Fetching traderoutes from rarethief...")
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, "", "")
	if err != nil {
		return Traderoutes{}, err
	}

	re, err := regexp.Compile(`var trade_routes\s*=\s*({.*})`)
	if err != nil {
		return Traderoutes{}, err
	}
	validJson := re.FindStringSubmatch(string(httpResp))

	var routes Traderoutes
	if len(validJson) > 1 {
		if err := json.Unmarshal([]byte(validJson[1]), &routes); err != nil {
			l.Errorf("Failed to unmarshal API response: %v", err)
			return Traderoutes{}, err
		}
		return routes, nil
	}

	return Traderoutes{}, fmt.Errorf("No traderoutes found.")
}
