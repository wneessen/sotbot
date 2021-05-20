package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/sotapi"
	"regexp"
	"strings"
)

// Just a test handler
func (b *Bot) SotReputation(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.SotReputation",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!rep") {
		l.Debugf("Received '!rep' request from user %v", m.Author.Username)
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

		wrongFormatMsg := fmt.Sprintf("Incorrect request format. Usage: !reputation " +
			"<athena|bilge|hoarder|hunter|merchant|order|reaper|seadog>")
		msgArray := strings.SplitN(m.Message.Content, " ", 2)
		if len(msgArray) != 2 {
			AnswerUser(s, m, wrongFormatMsg, m.Author.Mention())
			return
		}
		var validFaction = regexp.MustCompile(`^(athena|bilge|hoarder|hunter|merchant|order|reaper|seadog)$`)
		if !validFaction.MatchString(msgArray[1]) {
			AnswerUser(s, m, wrongFormatMsg, m.Author.Mention())
			return
		}
		validFactionMatch := validFaction.FindStringSubmatch(msgArray[1])
		if len(validFactionMatch) < 1 {
			AnswerUser(s, m, wrongFormatMsg, m.Author.Mention())
			return
		}

		userRatCookie := database.UserGetPrefString(b.Db, userObj.ID, "rat_cookie")
		if userRatCookie == "" {
			replyMsg := fmt.Sprintf("Sorry but you have no RAT cookie set. Try !setrat in the DMs")
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		userReputation, err := sotapi.GetFactionReputation(b.HttpClient, userRatCookie, strings.ToLower(validFactionMatch[0]))
		if err != nil {
			l.Errorf("An error occured fetching user progress: %v", err)
			replyMsg := fmt.Sprintf("Sorry but there was an error fetching your faction reputation"+
				" from the SoT API: %v", err)
			AnswerUser(s, m, replyMsg, m.Author.Mention())
			return
		}
		userMsg := fmt.Sprintf("You current reputation level with the %v faction (%q) is: %d (Rank: %v). "+
			"Your current level XP is %d. To reach the next reputation level (level %d), you need to reach a total"+
			" of %d XP. That's %d to go!",
			userReputation.Name, userReputation.Motto, userReputation.Level, userReputation.Rank,
			userReputation.Xp, userReputation.NextLevel.Level, userReputation.NextLevel.XpRequired,
			(userReputation.NextLevel.XpRequired - userReputation.Xp))
		AnswerUser(s, m, userMsg, m.Author.Mention())
	}
}
