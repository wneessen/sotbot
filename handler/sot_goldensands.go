package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
)

// GetSotGoldenSands fetches the current status around the battle of Golden Sands
func GetSotGoldenSands(h *http.Client, u *user.User) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotGoldenSands",
	})

	adventures, err := api.GetAdventures(h, u.RatCookie)
	if err != nil {
		l.Errorf("Adventures lookup failed: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	goals := adventures[0].GroupGoalProgress
	if goals == nil {
		return &discordgo.MessageEmbed{}, fmt.Errorf("the current adventure has no goals")
	}
	desc := "The situation is currently neutral. Neither the Reaper's nor Merrick's team is in the lead."
	if goals.LeadingGroupGoalID == "1463b6ce-7e0e-4b10-9765-ccff3d9a8152" {
		desc = "The force for the Flame is strong. The Reaper's are currently in the lead and on their way to burn down Golden Sands Outpost"
	}
	if goals.LeadingGroupGoalID == "f08b4609-0949-46c5-8fc9-3202d743474a" {
		desc = "Merrick's builder alliance is strong. Currently they are in the lead and on their way to rebuild Golden Sands Outpost"
	}
	responseEmbed := discordgo.MessageEmbed{
		Title:       "The battle for Golden Sands Outpost",
		Description: desc,
	}
	return &responseEmbed, nil
}
