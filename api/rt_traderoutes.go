package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
)

type Route struct {
	Outpost string `json:"outpost"`
	Sought  string `json:"sought_after"`
	Surplus string `json:"surplus"`
}

type Traderoutes struct {
	Dates  string           `json:"trade_route_dates"`
	Routes map[string]Route `json:"routes"`
}

func GetTraderoutes() (Traderoutes, error) {
	l := log.WithFields(log.Fields{
		"action": "rarethief.Traderoutes",
	})
	apiUrl := "https://maps.seaofthieves.rarethief.com/js/trade_routes.js"

	l.Debugf("Fetching traderoutes from rarethief...")
	httpResp, err := http.Get(apiUrl)
	if err != nil {
		return Traderoutes{}, err
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return Traderoutes{}, err
	}
	re, err := regexp.Compile(`var trade_routes\s*=\s*({.*})`)
	if err != nil {
		return Traderoutes{}, err
	}
	validJson := re.FindStringSubmatch(string(body))

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
