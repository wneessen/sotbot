package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/handler"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
	"strings"
)

// Let's the bot tell you the current date/time when requested via !time command
func (b *Bot) CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "bot.CommandHandler",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, "!") || len(m.Content) <= 1 {
		return
	}

	msgArray := strings.Split(m.Content, " ")
	command := strings.ToLower(msgArray[0])
	cmdNum := len(msgArray)

	var chanInfo *discordgo.Channel
	var err error
	chanInfo, err = s.Channel(m.ChannelID)
	if err != nil {
		l.Errorf("Failed to get channel info: %v", err)
	}

	userObj, err := user.NewUser(b.Db, b.Config, m.Author.ID)
	if err != nil {
		l.Errorf("Could not create user object: %v", err)
		return
	}
	userObj.AuthorName = m.Author.Username
	userObj.Mention = m.Author.Mention()

	if chanInfo.Type != discordgo.ChannelTypeDM {
		userPerm, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			l.Errorf("Failed to look up user-channel permission: %v", err)
		}
		userObj.ChanPermission = userPerm
	}

	switch {

	// SoT: Set RAT cookie
	case (command == "!setrat" || command == "!rat"):
		if !userObj.IsRegistered() {
			re := "Sorry, you are not registered"
			response.AnswerUser(s, m, re, true)
			return
		}
		if chanInfo.Type != discordgo.ChannelTypeDM {
			re := "You exposed your RAT cookie to a public channel. Please change your password immediately."
			response.AnswerUser(s, m, re, true)
			return
		}
		if cmdNum != 2 {
			re := "The !setrat command requires you to provide a cookie. Usage `!setrat <cookie>`"
			response.AnswerUser(s, m, re, true)
			return
		}
		re, err := handler.UserSetRatCookie(b.Db, b.Config, userObj, msgArray[1])
		if err != nil {
			re := fmt.Sprintf("An error occurred setting/updating your RAT cookie: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return
	}
}
