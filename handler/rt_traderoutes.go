package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"strings"
)

func GetTraderoutes() (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetTraderoutes",
	})
	traderoutes, err := api.GetTraderoutes()
	if err != nil {
		l.Errorf("An error occurred fetching traderoutes: %v", err)
		return "Sorry, couldn't fetch traderoutes", err
	}
	var response strings.Builder
	response.WriteString(fmt.Sprintf("Traderoutes for %v according to http://maps.seaofthieves.rarethief.com are\n\n", traderoutes.Dates))
	for _, v := range traderoutes.Routes {
		response.WriteString(fmt.Sprintf("**%v**: sought after ***%v*** / surplus ***%v***\n", v.Outpost, v.Sought, v.Surplus))
	}
	response.WriteString("\nSafe travels, don't let the reapers bite you...")

	return response.String(), nil
}
