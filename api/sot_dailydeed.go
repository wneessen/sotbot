package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
	"regexp"
)

type EventApiResponse struct {
	GlobalNav GlobalNav `json:"globalNav"`
}

type GlobalNav struct {
	Items []GlobalNavItem `json:"Items"`
}

type GlobalNavItem struct {
	Children []GlobalNavItemChild `json:"Children"`
}

type GlobalNavItemChild struct {
	DailyDeed DailyDeed `json:"DailyDeed"`
}

type DailyDeed struct {
	Title     string       `json:"Title"`
	Copy      string       `json:"Copy"`
	StartDate TimeRFC3339  `json:"StartDate"`
	EndDate   TimeRFC3339  `json:"EndDate"`
	Image     DailyDeedImg `json:"Image"`
}

type DailyDeedImg struct {
	Desktop string `json:"desktop"`
}

func GetDailyDeed(hc *http.Client, rc string) (DailyDeed, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetDailyDeed",
	})
	apiUrl := "https://www.seaofthieves.com/event-hub"

	l.Debugf("Fetching event-hub data from www.seaofthieves.com...")
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, rc, "https://www.seaofthieves.com/season-two")
	if err != nil {
		return DailyDeed{}, err
	}
	re := regexp.MustCompile(`<script>var APP_PROPS\s*=\s*({.*});</script>`)
	validJson := re.FindStringSubmatch(string(httpResp))
	if len(validJson) > 1 {
		var apiResponse EventApiResponse
		if err := json.Unmarshal([]byte(validJson[1]), &apiResponse); err != nil {
			l.Errorf("Failed to unmarshal API response: %v", err)
			return DailyDeed{}, err
		}
		for _, curItem := range apiResponse.GlobalNav.Items {
			if curItem.Children[0].DailyDeed.Title != "" {
				return curItem.Children[0].DailyDeed, nil
			}
		}
	}

	return DailyDeed{}, fmt.Errorf("No daily deed found in API response.")
}
