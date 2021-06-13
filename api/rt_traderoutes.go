package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Traderoutes defines the structure for the rarethief.com trading route
// API response
type Traderoutes struct {
	Dates     string           `json:"trade_route_dates"`
	Routes    map[string]Route `json:"routes"`
	ValidFrom time.Time
	ValidThru time.Time
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
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, nil, nil, true)
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
		dateArray := strings.SplitN(routes.Dates, " - ", 2)
		fromTime, err := time.Parse("2006/01/02", fmt.Sprintf("%v/%v", time.Now().Year(), dateArray[0]))
		if err == nil {
			routes.ValidFrom = fromTime
		}
		toTime, err := time.Parse("2006/01/02 15:04:05",
			fmt.Sprintf("%v/%v 23:59:59", time.Now().Year(), dateArray[1]))
		if err == nil {
			routes.ValidThru = toTime
		}
		return routes, nil
	}

	return Traderoutes{}, fmt.Errorf("No traderoutes found.")
}
