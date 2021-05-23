package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/httpclient"
	"github.com/wneessen/sotbot/random"
	"net/http"
)

func GetTMDbMovie(hc *http.Client, u string) (TMDbMovie, error) {
	l := log.WithFields(log.Fields{
		"action": "api.TMDbMovie",
	})

	var tmdbResp TMDbMovieResp
	l.Debugf("Fetching movie fact from TMDB API...")
	apiUrl := fmt.Sprintf("https://api.themoviedb.org%v", u)
	httpResp, err := httpclient.HttpReqGet(apiUrl, hc, "", "")
	if err != nil {
		return TMDbMovie{}, err
	}
	if err := json.Unmarshal(httpResp, &tmdbResp); err != nil {
		l.Errorf("Failed to unmarshal API response: %v", err)
		return TMDbMovie{}, err
	}
	numResults := len(tmdbResp.Results)
	randNum, err := random.Number(numResults)
	if err != nil {
		return TMDbMovie{}, err
	}

	return tmdbResp.Results[randNum], nil
}
