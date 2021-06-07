package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/user"
	"net/http"
	"regexp"
	"strings"
)

// Just a test handler
func GetSotReputation(h *http.Client, u *user.User, f string) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetSotReputation",
	})

	wrongFormatMsg := fmt.Sprintf("Incorrect request format. Usage: !reputation " +
		"<athena|bilge|hoarder|hunter|merchant|order|reaper|seadog>")
	validFaction, err := regexp.Compile(`^(athena|bilge|hoarder|hunter|merchant|order|reaper|seadog)$`)
	if err != nil {
		return "", err
	}
	if !validFaction.MatchString(f) {
		return wrongFormatMsg, nil
	}
	validFactionMatch := validFaction.FindStringSubmatch(f)
	if len(validFactionMatch) < 1 {
		return wrongFormatMsg, nil
	}
	userReputation, err := api.GetFactionReputation(h, u.RatCookie, strings.ToLower(validFactionMatch[0]))
	if err != nil {
		l.Errorf("An error occurred fetching user progress: %v", err)
		return "", err
	}
	responseMsg := fmt.Sprintf("You current reputation level with the **%v faction** (%q) is: **%d (Rank: %v)**. "+
		"Your current level **XP is %d**. To reach the next reputation level (level %d), you need to reach a total"+
		" of **%d XP**. That's **%d to go**!",
		userReputation.Name, userReputation.Motto, userReputation.Level, userReputation.Rank,
		userReputation.Xp, userReputation.NextLevel.Level, userReputation.NextLevel.XpRequired,
		userReputation.NextLevel.XpRequired-userReputation.Xp)
	return responseMsg, nil
}
