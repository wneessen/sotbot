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

type BandTitle struct {
	AthenasFortune   []string
	GoldHoarder      []string
	MerchantAlliance []string
	OrderOfSouls     []string
	ReapersBone      []string
}

type FactionLedger struct {
	Name       string
	Band       int `json:"band"`
	BandTitle  string
	Rank       int `json:"rank"`
	Score      int `json:"score"`
	ToNextRank int `json:"toNextRank"`
}

func GetFactionLedger(hc *http.Client, rc string, f string) (FactionLedger, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetFactionLedger",
	})

	factionBandTitle := BandTitle{
		AthenasFortune:   []string{"Legend", "Guardian", "Voyager", "Seeker"},
		GoldHoarder:      []string{"Captain", "Marauder", "Seafarer", "Castaway"},
		MerchantAlliance: []string{"Admiral", "Commander", "Cadet", "Sailor"},
		OrderOfSouls:     []string{"Grandee", "Chief", "Mercenary", "Apprentice"},
		ReapersBone:      []string{"Master", "Keeper", "Servant", "Follower"},
	}

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
		userFactionLedger.BandTitle = factionBandTitle.AthenasFortune[userFactionLedger.Band]
	case "hoarder":
		userFactionLedger.Name = "Gold Hoarders"
		userFactionLedger.BandTitle = factionBandTitle.GoldHoarder[userFactionLedger.Band]
	case "merchant":
		userFactionLedger.Name = "Merchant Alliance"
		userFactionLedger.BandTitle = factionBandTitle.MerchantAlliance[userFactionLedger.Band]
	case "order":
		userFactionLedger.Name = "Order of Souls"
		userFactionLedger.BandTitle = factionBandTitle.OrderOfSouls[userFactionLedger.Band]
	case "reaper":
		userFactionLedger.Name = "Reaper's Bones"
		userFactionLedger.BandTitle = factionBandTitle.ReapersBone[userFactionLedger.Band]
	}

	return userFactionLedger, nil
}
