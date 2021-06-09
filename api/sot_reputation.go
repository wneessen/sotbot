package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type Factions struct {
	AthenasFortune   FactionReputation `json:"AthenasFortune"`
	BilgeRats        FactionReputation `json:"BilgeRats"`
	GoldHoarders     FactionReputation `json:"GoldHoarders"`
	HuntersCall      FactionReputation `json:"HuntersCall"`
	MerchantAlliance FactionReputation `json:"MerchantAlliance"`
	OrderOfSouls     FactionReputation `json:"OrderOfSouls"`
	ReapersBones     FactionReputation `json:"ReapersBones"`
	SeaDogs          FactionReputation `json:"SeaDogs"`
}

type FactionReputation struct {
	Name      string
	Motto     string           `json:"Motto"`
	Level     int              `json:"Level"`
	Rank      string           `json:"Rank"`
	Xp        int              `json:"Xp"`
	NextLevel FactionNextLevel `json:"NextCompanyLevel"`
}

type FactionNextLevel struct {
	Level      int `json:"Level"`
	XpRequired int `json:"XpRequiredToAttain"`
}

func GetFactionReputation(hc *http.Client, rc string, f string) (FactionReputation, error) {
	l := log.WithFields(log.Fields{
		"action": "sotapi.GetFactionReputation",
	})
	apiUrl := "https://www.seaofthieves.com/api/profilev2/reputation"

	l.Debugf("Fetching user reputation in %v faction from API...", f)
	var userReps Factions
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, &rc, nil, false)
	if err != nil {
		return FactionReputation{}, err
	}
	if err := json.Unmarshal(httpResp, &userReps); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return FactionReputation{}, err
	}

	switch f {
	case "athena":
		userReps.AthenasFortune.Name = "Athenas Fortune"
		return userReps.AthenasFortune, nil
	case "bilge":
		userReps.BilgeRats.Name = "Bilge Rats"
		return userReps.BilgeRats, nil
	case "hoarder":
		userReps.GoldHoarders.Name = "Gold Hoarders"
		return userReps.GoldHoarders, nil
	case "hunter":
		userReps.HuntersCall.Name = "Hunter's Call"
		return userReps.HuntersCall, nil
	case "merchant":
		userReps.MerchantAlliance.Name = "Merchant Alliance"
		return userReps.MerchantAlliance, nil
	case "order":
		userReps.OrderOfSouls.Name = "Order of Souls"
		return userReps.OrderOfSouls, nil
	case "reaper":
		userReps.ReapersBones.Name = "Reaper's Bones"
		return userReps.ReapersBones, nil
	case "seadog":
		userReps.SeaDogs.Name = "Sea Dogs"
		return userReps.SeaDogs, nil
	default:
		l.Errorf("Wrong faction name provided")
		return FactionReputation{}, fmt.Errorf("Unknown faction")
	}
}
