package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/sotapi"
)

// Just a test handler
func (b *Bot) LatestAchievement(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.TestHandler",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.Content == "!achievement" {
		l.Debugf("Received '!achievement' request from user %v", m.Author.Username)
		userObj, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			l.Errorf("Database lookup failed: %v", err)
			return
		}
		if userObj.ID <= 0 {
			replyMsg := fmt.Sprintf("%v, sorry but your are not a registered user.",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
			return
		}
		userRatCookie := database.UserGetPrefString(b.Db, userObj.ID, "rat_cookie")
		if userRatCookie == "" {
			replyMsg := fmt.Sprintf("%v, sorry but you have no RAT cookie set. Try !setrat in the DMs",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
			return
		}
		userAchievement, err := sotapi.GetLatestAchievement(b.HttpClient, userRatCookie)
		if err != nil {
			l.Errorf("An error occured fetching user achievements: %v", err)
			replyMsg := fmt.Sprintf("Sorry, %v but there was an error fetching your achievements"+
				" from the SoT API: %v", m.Author.Mention(), err)
			AnswerUser(s, m, replyMsg)
			return
		}
		embedTitle := fmt.Sprintf("%v, your latest achievement is: %v",
			m.Author.Username, userAchievement.Name)
		messageEmbed := discordgo.MessageEmbed{
			Title:       embedTitle,
			Description: userAchievement.Description,
			Image: &discordgo.MessageEmbedImage{
				URL: userAchievement.MediaUrl,
			},
			Type: discordgo.EmbedTypeImage,
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &messageEmbed)
		if err != nil {
			l.Errorf("Failed to send embeded message: %v", err)
			return
		}
	}
}
