package handler

import (
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"net/http"
)

func RandomFact(h *http.Client) (string, bool, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.RandomFact",
	})

	randFact, err := api.GetRandFact(h)
	if err != nil {
		l.Errorf("Could not fetch random fact: %v", err)
		return "", false, err
	}

	return randFact.Text, true, nil
}
