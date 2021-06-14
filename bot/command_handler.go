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

	// Placeholder for the legacy !airhorn command
	case command == "!airhorn" && cmdNum == 1:
		re := "The !airhorn command is now a slash-command. Try `/play airhorn` instead"
		response.AnswerUser(s, m, re, true)
		return

	// Tell us the current time
	case command == "!time" && cmdNum == 1:
		re := handler.TellTime()
		response.AnswerUser(s, m, re, true)
		return

	// Show some memory statistics
	case (command == "!memory" || command == "!mem") && cmdNum == 1:
		if !userObj.IsAdmin() {
			return
		}
		re := handler.TellMemUsage()
		response.AnswerUser(s, m, re, false)
		return

	// Show bot's uptime
	case (command == "!uptime" || command == "!up") && cmdNum == 1:
		re, err := handler.Uptime(b.StartTime)
		if err != nil {
			re = fmt.Sprintf("Sorry, an error occurred calculating the uptime: %v", err)
		}
		response.AnswerUser(s, m, re, true)
		return

	// Reply with random useless fact
	case command == "!fact" && cmdNum == 1:
		re, err := handler.RandomFact(b.HttpClient)
		if err != nil {
			re = fmt.Sprintf("An error occurred while fetching the random fact API: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return

	// Check Urban dictionary
	case command == "!ud" || command == "!urban":
		if cmdNum == 1 {
			em, err := handler.UrbanDict(b.HttpClient, "")
			if err != nil {
				re := fmt.Sprintf("An error occurred while fetching the urban dictionary API: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
			response.Embed(s, chanInfo.ID, em)
			return
		}
		if cmdNum == 2 {
			em, err := handler.UrbanDict(b.HttpClient, msgArray[1])
			if err != nil {
				re := fmt.Sprintf("An error occurred while fetching the urban dictionary API: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
			response.Embed(s, chanInfo.ID, em)
			return
		}

	// Get a random movie recommendation
	case command == "!movie":
		if cmdNum == 1 {
			if b.TMDb == nil {
				re := "You haven't specified a TMDb API key in your config file."
				response.AnswerUser(s, m, re, true)
				return
			}
			em, err := handler.TMDbRandMovie(b.TMDb)
			if err != nil {
				re := fmt.Sprintf("An error occurred while fetching the TMDB API: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
			response.Embed(s, chanInfo.ID, em)
			return
		}
		if cmdNum > 1 {
			if b.TMDb == nil {
				re := "You haven't specified a TMDb API key in your config file."
				response.AnswerUser(s, m, re, true)
				return
			}
			em, err := handler.TMDbSearchMovie(b.TMDb, msgArray[1:])
			if err != nil {
				if err.Error() == "No matching movie found" {
					re := "Sorry, but I wasn't able to find a movie matching your search criteria."
					response.AnswerUser(s, m, re, true)
					return
				}
				re := fmt.Sprintf("An error occurred while fetching the TMDB API: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
			response.Embed(s, chanInfo.ID, em)
			return
		}

	// Get a random movie recommendation
	case command == "!tv":
		if cmdNum == 1 {
			if b.TMDb == nil {
				re := "You haven't specified a TMDb API key in your config file."
				response.AnswerUser(s, m, re, true)
				return
			}
			em, err := handler.TMDbRandTvSeries(b.TMDb)
			if err != nil {
				re := fmt.Sprintf("An error occurred while fetching the TMDB API: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
			response.Embed(s, chanInfo.ID, em)
			return
		}
		if cmdNum > 1 {
			if b.TMDb == nil {
				re := "You haven't specified a TMDb API key in your config file."
				response.AnswerUser(s, m, re, true)
				return
			}
			em, err := handler.TMDbSearchTvSeries(b.TMDb, msgArray[1:])
			if err != nil {
				if err.Error() == "No matching TV series found" {
					re := "Sorry, but I wasn't able to find a TV series matching your search criteria."
					response.AnswerUser(s, m, re, true)
					return
				}
				re := fmt.Sprintf("An error occurred while fetching the TMDB API: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
			response.Embed(s, chanInfo.ID, em)
			return
		}

	// SoT: Show user's balance
	case (command == "!balance" || command == "!bal") && cmdNum == 1:
		if !userObj.IsRegistered() {
			return
		}
		if !userObj.HasRatCookie() {
			return
		}
		re, err := handler.GetSotBalance(b.Db, b.HttpClient, userObj)
		if err != nil {
			if err.Error() == "notify" {
				dmMsg := fmt.Sprintf("The last 3 attempts to communicate with the SoT API failed. " +
					"This likely means, that your RAT cookie has expired. Please use the !setrat function to " +
					"update your cookie.")
				response.DmUser(s, userObj, dmMsg, true, false)
			} else {
				re = fmt.Sprintf("An error occurred checking your SoT balance: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
		}
		response.AnswerUser(s, m, re, true)
		return

	// SoT: Show user's season progress
	case command == "!season" && cmdNum == 1:
		if !userObj.IsRegistered() {
			return
		}
		if !userObj.HasRatCookie() {
			return
		}
		re, err := handler.GetSotSeasonProgress(b.HttpClient, userObj)
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT season progress: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return

	// SoT: Show user's reputation with a faction/company
	case (command == "!reputation" || command == "!rep"):
		if !userObj.IsRegistered() {
			return
		}
		if !userObj.HasRatCookie() {
			return
		}
		var re string
		var err error
		if cmdNum != 2 {
			re, err = handler.GetSotReputation(b.HttpClient, userObj, "")
		}
		if cmdNum == 2 {
			re, err = handler.GetSotReputation(b.HttpClient, userObj, msgArray[1])
		}
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT reputation level: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return

	// SoT: Show user's ledger position with a faction/company
	case (command == "!ledger" || command == "!led"):
		if !userObj.IsRegistered() {
			return
		}
		if !userObj.HasRatCookie() {
			return
		}
		var re string
		var err error
		if cmdNum != 2 {
			re, err = handler.GetSotLedger(b.HttpClient, userObj, "")
		}
		if cmdNum == 2 {
			re, err = handler.GetSotLedger(b.HttpClient, userObj, msgArray[1])
		}
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT ledger rank: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return

	// SoT: Show user's general stats
	case (command == "!stats" || command == "!stat") && cmdNum == 1:
		if !userObj.IsRegistered() {
			return
		}
		if !userObj.HasRatCookie() {
			return
		}
		re, err := handler.GetSotStats(b.HttpClient, userObj)
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT general stats: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.AnswerUser(s, m, re, true)
		return

	// SoT: Quote a random SoT pirate code article
	case command == "!code" && cmdNum == 1:
		em, err := handler.GetSotRandomCode()
		if err != nil {
			re := fmt.Sprintf("An error occurred quoting the SoT pirate code: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.Embed(s, m.ChannelID, em)
		return

	// SoT: Set RAT cookie
	case (command == "!setrat" || command == "!rat"):
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

	// RareThief: Get Traderoutes
	case command == "!tr" && cmdNum == 1:
		em, err := handler.GetTraderoutes(b.HttpClient, b.Db)
		if err != nil {
			re := fmt.Sprintf("An error occurred fetching traderoutes: %v", err)
			response.AnswerUser(s, m, re, true)
			return
		}
		response.Embed(s, m.ChannelID, em)
		return

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

	// OWM: Current weather
	case command == "!weather":
		if cmdNum == 1 {
			re := "You need to tell me a location as well. Usage: `!weather <location>`"
			response.AnswerUser(s, m, re, true)
			return
		}
		if cmdNum > 1 {
			if b.OwmClient == nil {
				re := "You haven't specified a OpenWeatherMap API key in your config file."
				response.AnswerUser(s, m, re, true)
				return
			}
			re, err := handler.GetCurrentWeather(b.OwmClient, msgArray[1:])
			if err != nil {
				re := fmt.Sprintf("An error occurred fetching OWM weather data: %v", err)
				response.AnswerUser(s, m, re, true)
				return
			}
			response.AnswerUser(s, m, re, true)
			return
		}
	}
}
