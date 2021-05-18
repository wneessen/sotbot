package factapi

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"net/http"
)

type RandomFact struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	Source    string `json:"source"`
	SourceUrl string `json:"source_url"`
	Language  string `json:"language"`
	Permalink string `json:"permalink"`
}

func GetRandFact(hc *http.Client) (RandomFact, error) {
	l := log.WithFields(log.Fields{
		"action": "factapi.GetRandFact",
	})
	var randomFact RandomFact
	apiUrl := "https://uselessfacts.jsph.pl/random.json?language=en"

	l.Debugf("Fetching random fact from API...")
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, "", "")
	if err != nil {
		return randomFact, err
	}
	if err := json.Unmarshal(httpResp, &randomFact); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return randomFact, err
	}

	return randomFact, nil
}
