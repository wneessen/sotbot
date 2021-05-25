package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type UrbanDictEntry struct {
	Definition string `json:"definition"`
	PermaLink  string `json:"permalink"`
	Example    string `json:"example"`
	Author     string `json:"author"`
	Word       string `json:"word"`
}

type UrbanDict struct {
	List []UrbanDictEntry `json:"list"`
}

func GetUrbanDict(hc *http.Client, w string) (UrbanDictEntry, error) {
	l := log.WithFields(log.Fields{
		"action": "api.GetUrbanDict",
	})

	var urbanReps UrbanDict
	var apiUrl string
	if w == "" {
		apiUrl = "https://api.urbandictionary.com/v0/random"
	}
	if w != "" {
		apiUrl = fmt.Sprintf("https://api.urbandictionary.com/v0/define?term=%v", w)
	}

	l.Debugf("Fetching UD definition fact from UD API...")
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, "", "")
	if err != nil {
		return UrbanDictEntry{}, err
	}
	if err := json.Unmarshal(httpResp, &urbanReps); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return UrbanDictEntry{}, err
	}

	return urbanReps.List[0], nil
}
