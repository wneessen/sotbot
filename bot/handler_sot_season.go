package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/sotapi"
)

// Just a test handler
func (b *Bot) SeasonProgress(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.TestHandler",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.Content == "!season" {
		l.Debugf("Received '!season' request from user %v", m.Author.Username)
		userObj, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			l.Errorf("Database lookup failed: %v", err)
			return
		}
		if userObj.ID <= 0 {
			replyMsg := fmt.Sprintf("Sorry but your are not a registered user.")
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		userRatCookie := database.UserGetPrefString(b.Db, userObj.ID, "rat_cookie")
		if userRatCookie == "" {
			replyMsg := fmt.Sprintf("Sorry but you have no RAT cookie set. Try !setrat in the DMs")
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		userAchievement, err := sotapi.GetSeasonProgress(b.HttpClient, userRatCookie)
		if err != nil {
			l.Errorf("An error occured fetching user progress: %v", err)
			replyMsg := fmt.Sprintf("Sorry but there was an error fetching your season progress"+
				" from the SoT API: %v", err)
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		userMsg := fmt.Sprintf("You are currently sailing in %v. Your renown level is %d (Tier: %d). "+
			"Of the total amount of %d season challanges, so far, you completed %d.", userAchievement.SeasonTitle,
			userAchievement.LevelProgress, userAchievement.Tier, userAchievement.TotalChallenges,
			userAchievement.CompletedChallenges)
		AnswerUser(s, m, userMsg, m.Author.Mention())
	}
}
