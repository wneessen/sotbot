package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type ApiLedger struct {
	Current CurrentLedger `json:"current"`
}

type CurrentLedger struct {
	Friends FriendsLedger `json:"friends"`
}

type FriendsLedger struct {
	User FactionLedger `json:"user"`
}

type FactionLedger struct {
	Name       string
	Band       int `json:"band"`
	Rank       int `json:"rank"`
	Score      int `json:"score"`
	ToNextRank int `json:"toNextRank"`
}

func GetFactionLedger(hc *http.Client, rc string, f string) (FactionLedger, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetFactionLedger",
	})

	var apiUrl string
	apiBase := "https://www.seaofthieves.com/api/ledger/friends/"

	switch f {
	case "athena":
		apiUrl = fmt.Sprintf("%v%v?count=10", apiBase, "AthenasFortune")
	case "hoarder":
		apiUrl = fmt.Sprintf("%v%v?count=10", apiBase, "GoldHoarders")
	case "merchant":
		apiUrl = fmt.Sprintf("%v%v?count=10", apiBase, "MerchantAlliance")
	case "order":
		apiUrl = fmt.Sprintf("%v%v?count=10", apiBase, "OrderOfSouls")
	case "reaper":
		apiUrl = fmt.Sprintf("%v%v?count=10", apiBase, "ReapersBones")
	default:
		return FactionLedger{}, fmt.Errorf("Unknown faction")
	}

	l.Debugf("Fetching user ledger position in %v faction from API...", f)
	var userApiLedger ApiLedger
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, rc, "")
	if err != nil {
		return FactionLedger{}, err
	}
	if err := json.Unmarshal(httpResp, &userApiLedger); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return FactionLedger{}, err
	}

	userFactionLedger := userApiLedger.Current.Friends.User
	switch f {
	case "athena":
		userFactionLedger.Name = "Athenas Fortune"
	case "hoarder":
		userFactionLedger.Name = "Gold Hoarders"
	case "merchant":
		userFactionLedger.Name = "Merchant Alliance"
	case "order":
		userFactionLedger.Name = "Order of Souls"
	case "reaper":
		userFactionLedger.Name = "Reaper's Bones"
	}

	return userFactionLedger, nil
}
