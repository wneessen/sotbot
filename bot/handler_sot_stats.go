package bot

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database"
	"github.com/wneessen/sotbot/httpclient"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
)

type ApiBalance struct {
	GamerTag     string `json:"gamertag"`
	Title        string `json:"title"`
	Doubloons    int    `json:"doubloons"`
	Gold         int    `json:"gold"`
	AncientCoins int    `json:"ancientCoins"`
}

// Set a SoT RAT cookie
func (b *Bot) SetRatCookie(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.SetRatCookie",
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

	if strings.HasPrefix(m.Message.Content, "!setrat") {
		l.Debugf("Received '!setrat' request from user %v", m.Author.Username)
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

		wrongFormatMsg := fmt.Sprintf("%v, incorrect request format. Usage: !setrat <ratcookie>",
			m.Author.Mention())
		msgArray := strings.SplitN(m.Message.Content, " ", 2)
		if len(msgArray) != 2 {
			AnswerUser(s, m, wrongFormatMsg)
			return
		}

		if err := database.UserSetPref(b.Db, userObj.ID, "rat_cookie", msgArray[1]); err != nil {
			l.Errorf("Failed to store RAT cookie in DB: %v", err)
			replyMsg := fmt.Sprintf("Sorry, I couldn't store/update your cookie in the DB.")
			AnswerUser(s, m, replyMsg)
			return
		}
		replyMsg := fmt.Sprintf("%v, thanks for setting/updating your RAT cookie.",
			m.Author.Mention())
		AnswerUser(s, m, replyMsg)
	}
}

// Get current SoT balance
func (b *Bot) GetBalance(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetBalance",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!balance") {
		l.Debugf("Received '!balance' request from user %v", m.Author.Username)
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

		apiUrl := "https://www.seaofthieves.com/api/profilev2/balance"
		l.Debugf("Fetching balance from API...")
		httpResp, err := httpclient.HttpReqGet(apiUrl, b.HttpClient, userRatCookie, "")
		if err != nil {
			l.Errorf("Failed to fetch balance from API: %v", err)
			replyMsg := fmt.Sprintf("%v, sorry an error occured while fetching your SoT balance from API.",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
			return
		}
		var userBalance ApiBalance
		if err := json.Unmarshal(httpResp, &userBalance); err != nil {
			l.Errorf("Failed to unmarshal API response: %v", err)
			replyMsg := fmt.Sprintf("%v, sorry I wasn't able to understand the API response JSON.",
				m.Author.Mention())
			AnswerUser(s, m, replyMsg)
			return
		}

		p := message.NewPrinter(language.German)
		replyMsg := fmt.Sprintf("%v, your current SoT balance is: %v gold, %v doubloons and %v ancient coins",
			m.Author.Mention(),
			p.Sprintf("%d", userBalance.Gold),
			p.Sprintf("%d", userBalance.Doubloons),
			p.Sprintf("%d", userBalance.AncientCoins),
		)
		AnswerUser(s, m, replyMsg)
	}
}
