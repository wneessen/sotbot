package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/user"
	"regexp"
	"strings"
)

// Self-check if a user is registered
func (b *Bot) CurrentUserIsRegistered(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.CurrentUserIsRegistered",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	curChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		l.Errorf("Failed to get channel info: %v", err)
		return
	}
	if curChannel.Type != discordgo.ChannelTypeDM {
		return
	}

	if m.Message.Content == "!userinfo" {
		l.Debugf("Received '!userinfo' request from user %v", m.Author.Username)
		userObj, err := database.GetUser(b.Db, m.Author.ID)
		if err != nil {
			l.Errorf("Database lookup failed: %v", err)
			return
		}

		var returnMsg string
		if userObj.ID <= 0 {
			returnMsg = "Sorry, you are not a registered user."
			AnswerUser(s, m, returnMsg)
			return
		}

		returnMsg = "You are a registered user."
		AnswerUser(s, m, returnMsg)
	}
}

// Register a new user
func (b *Bot) RegisterUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.RegisterUser",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	curChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		l.Errorf("Failed to get channel info: %v", err)
		return
	}
	if curChannel.Type != discordgo.ChannelTypeGuildText {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!register") {
		l.Debugf("Received '!register' request from user %v", m.Author.Username)

		if !user.IsAdmin(s, m.Author.ID, m.ChannelID) {
			l.Debugf("User is not an admin user")
			return
		}

		wrongFormatMsg := fmt.Sprintf("%v, incorrect request format. Usage: !register <@user>",
			m.Author.Mention())
		msgArray := strings.SplitN(m.Message.Content, " ", 2)
		if len(msgArray) != 2 {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}
		var validUser = regexp.MustCompile(`^\<@[\!&]*(\d+)>$`)
		if !validUser.MatchString(msgArray[1]) {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}
		validUserMatches := validUser.FindStringSubmatch(msgArray[1])
		if len(validUserMatches) < 2 {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}
		dbUser, err := database.GetUser(b.Db, validUserMatches[1])
		if err != nil {
			l.Errorf("Failed to look up user in database: %v", err)
			replyMsg := fmt.Sprintf("%v, unfortunately I was not able to store the user in the database",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
		}
		if dbUser.ID > 0 {
			replyMsg := fmt.Sprintf("%v, user %v is already registered.",
				m.Author.Mention(), validUserMatches[0])
			AnswerUser(s, m, replyMsg)
			return
		}

		if err := database.CreateUser(b.Db, validUserMatches[1]); err != nil {
			l.Errorf("Failed to store user in database: %v", err)
			replyMsg := fmt.Sprintf("%v, unfortunately I was not able to store the user in the database",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)

		}

		replyMsg := fmt.Sprintf("%v, user %v successfully registered.",
			m.Author.Mention(), validUserMatches[0])
		AnswerUser(s, m, replyMsg)
	}
}

// UnRegister a user
func (b *Bot) UnRegisterUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.UnRegisterUser",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	curChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		l.Errorf("Failed to get channel info: %v", err)
		return
	}
	if curChannel.Type != discordgo.ChannelTypeGuildText {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!unregister") {
		l.Debugf("Received '!unregister' request from user %v", m.Author.Username)

		if !user.IsAdmin(s, m.Author.ID, m.ChannelID) {
			l.Debugf("User is not an admin user")
			return
		}

		wrongFormatMsg := fmt.Sprintf("%v, incorrect request format. Usage: !unregister <@user>",
			m.Author.Mention())
		msgArray := strings.SplitN(m.Message.Content, " ", 2)
		if len(msgArray) != 2 {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}
		var validUser = regexp.MustCompile(`^\<@[\!&]*(\d+)>$`)
		if !validUser.MatchString(msgArray[1]) {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}
		validUserMatches := validUser.FindStringSubmatch(msgArray[1])
		if len(validUserMatches) < 2 {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}
		dbUser, err := database.GetUser(b.Db, validUserMatches[1])
		if err != nil {
			l.Errorf("Failed to look up user in database: %v", err)
			replyMsg := fmt.Sprintf("%v, unfortunately I was not able to unregister the user.",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
		}
		if dbUser.ID <= 0 {
			replyMsg := fmt.Sprintf("%v, user %v is not registered.",
				m.Author.Mention(), validUserMatches[0])
			AnswerUser(s, m, replyMsg)
			return
		}

		if err := database.DeleteUser(b.Db, &dbUser); err != nil {
			l.Errorf("Failed to delete user in database: %v", err)
			replyMsg := fmt.Sprintf("%v, unfortunately I was not able to unregister the user.",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)

		}

		replyMsg := fmt.Sprintf("%v, user %v successfully unregistered.",
			m.Author.Mention(), validUserMatches[0])
		AnswerUser(s, m, replyMsg)
	}
}
