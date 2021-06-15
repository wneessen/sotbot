package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/handler"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
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

	l.Debugf("Received a %q command from %v", command, userObj.AuthorName)
	switch {

	// User management: Tell the user if they are registered
	case (command == "!userinfo" || command == "!info") && cmdNum == 1:
		re := handler.UserIsRegistered(userObj)
		response.AnswerUser(s, m, re, true)
		return

	// User management: Register a new user
	case (command == "!register" || command == "!reg"):
		if !userObj.IsAdmin() {
			return
		}
		if chanInfo.Type != discordgo.ChannelTypeGuildText {
			return
		}
		var re string
		var err error
		if cmdNum == 2 {
			re, err = handler.RegisterUser(b.Db, msgArray[1])
		}
		if cmdNum != 2 {
			re, err = handler.RegisterUser(b.Db, "")
		}
		if err != nil {
			re := fmt.Sprintf("An error occurred registering user: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return

	// User management: Un-register a new user
	case (command == "!unregister" || command == "!unreg"):
		if !userObj.IsAdmin() {
			return
		}
		if chanInfo.Type != discordgo.ChannelTypeGuildText {
			return
		}
		var re string
		var err error
		if cmdNum == 2 {
			re, err = handler.UnregisterUser(b.Db, msgArray[1])
		}
		if cmdNum != 2 {
			re, err = handler.UnregisterUser(b.Db, msgArray[1])
		}
		if err != nil {
			re := fmt.Sprintf("An error occurred registering user: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return
	}
}
