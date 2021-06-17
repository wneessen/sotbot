package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
	"regexp"
	"time"
)

type EventApiResponse struct {
	Data EventData `json:"data"`
}

type EventData struct {
	Components []EventDataComponent `json:"components"`
}

type EventDataComponent struct {
	Data EventDataComponentData `json:"data"`
}

type EventDataComponentData struct {
	BountyList []BountyListApi `json:"BountyList"`
}

type BountyListApi struct {
	Type                      string          `json:"#Type"`
	Title                     string          `json:"Title"`
	BodyText                  string          `json:"BodyText"`
	StartDateApi              *ApiTimeRFC3339 `json:"StartDate,omitempty"`
	EndDateApi                *ApiTimeRFC3339 `json:"EndDate,omitempty"`
	CompletedAtApi            *ApiTimeRFC3339 `json:"CompletedAt,omitempty"`
	Image                     DailyDeedImg    `json:"Image"`
	EntitlementRewardValue    int             `json:"EntitlementRewardValue"`
	EntitlementRewardCurrency string          `json:"EntitlementRewardCurrency"`
}

type BountyList struct {
	Type                      string
	Title                     string
	BodyText                  string
	StartDate                 time.Time
	EndDate                   time.Time
	CompletedAt               time.Time
	Image                     DailyDeedImg
	EntitlementRewardValue    int
	EntitlementRewardCurrency string
}

type DailyDeedImg struct {
	Desktop string `json:"desktop"`
}

func GetDailyDeed(hc *http.Client, rc string) (BountyList, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetDailyDeed",
	})
	apiUrl := "https://www.seaofthieves.com/event-hub"

	l.Debugf("Fetching event-hub data from www.seaofthieves.com...")
	ref := "https://www.seaofthieves.com/season-two"
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, &rc, &ref, false)
	if err != nil {
		return BountyList{}, err
	}
	re, err := regexp.Compile(`<script>var APP_PROPS\s*=\s*({.*});</script>`)
	if err != nil {
		return BountyList{}, err
	}
	validJson := re.FindStringSubmatch(string(httpResp))
	if len(validJson) > 1 {
		var apiResponse EventApiResponse
		if err := json.Unmarshal([]byte(validJson[1]), &apiResponse); err != nil {
			l.Errorf("Failed to unmarshal API response: %v", err)
			return BountyList{}, err
		}
		nowTime := time.Now().Unix()
		for _, curComp := range apiResponse.Data.Components {
			for _, curBounty := range curComp.Data.BountyList {
				if curBounty.StartDateApi.Time().Unix() <= nowTime && curBounty.EndDateApi.Time().Unix() >= nowTime {
					returnBounty := BountyList{
						Type:                      curBounty.Type,
						Title:                     curBounty.Title,
						BodyText:                  curBounty.BodyText,
						Image:                     curBounty.Image,
						EntitlementRewardValue:    curBounty.EntitlementRewardValue,
						EntitlementRewardCurrency: curBounty.EntitlementRewardCurrency,
					}
					if curBounty.StartDateApi != nil {
						returnBounty.StartDate = curBounty.StartDateApi.Time()
					}
					if curBounty.EndDateApi != nil {
						returnBounty.EndDate = curBounty.EndDateApi.Time()
					}
					if curBounty.CompletedAtApi != nil {
						returnBounty.CompletedAt = curBounty.CompletedAtApi.Time()
					}
					return returnBounty, nil
				}
			}
		}
	}

	return BountyList{}, fmt.Errorf("No daily deed found in API response.")
}
