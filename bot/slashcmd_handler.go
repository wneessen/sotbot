package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/handler"
	"github.com/wneessen/sotbot/response"
	"github.com/wneessen/sotbot/user"
)

// SlashCmdHandler handles all incoming slash commands. It creates a general user object
// and forward the command data to the corresponding handler submethod
func (b *Bot) SlashCmdHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := log.WithFields(log.Fields{
		"action": "bot.SlashCmdHandler",
	})

	cmdName := i.Data.Name

	var chanInfo *discordgo.Channel
	var err error
	chanInfo, err = s.Channel(i.ChannelID)
	if err != nil {
		l.Errorf("Failed to get channel info: %v", err)
	}

	var userObj *user.User
	if chanInfo != nil && chanInfo.Type != discordgo.ChannelTypeDM {
		userObj, err = user.NewUser(b.Db, b.Config, i.Member.User.ID)
		if err != nil {
			l.Errorf("Could not create user object: %v", err)
			return
		}
		userObj.AuthorName = i.Member.User.Username
		userObj.Mention = i.Member.User.Mention()

		userPerm, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			l.Errorf("Failed to look up user-channel permission: %v", err)
		}
		userObj.ChanPermission = userPerm
	}
	if chanInfo != nil && chanInfo.Type == discordgo.ChannelTypeDM {
		userObj, err = user.NewUser(b.Db, b.Config, i.User.ID)
		if err != nil {
			l.Errorf("Could not create user object: %v", err)
			return
		}
		userObj.AuthorName = i.User.Username
		userObj.Mention = i.User.Mention()
	}

	switch {
	// DM the commands help to the user
	case cmdName == "help":
		helpTexts := handler.Help()
		response.SlashCmdResponseMention(s, i.Interaction, userObj, "Please check your DMs.")
		for _, textPart := range helpTexts {
			response.DmUser(s, userObj, "`"+textPart+"`", false, true)
		}
		return

	// Reply with the version string of the bot
	case cmdName == "version":
		re := handler.TellVersion()
		response.SlashCmdResponseMention(s, i.Interaction, userObj, re)
		return

	// Tell the current time
	case cmdName == "time":
		re := handler.TellTime()
		response.SlashCmdResponseMention(s, i.Interaction, userObj, re)
		return

	// Tell the bot's uptime
	case cmdName == "uptime":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		re, err := handler.Uptime(b.StartTime)
		if err != nil {
			response.SlashCmdResponseEdit(s, i.Interaction, userObj,
				fmt.Sprintf("An error occured calculating the bots uptime: %v", err), true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// Tell the bot's memory usage
	case cmdName == "memory":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsAdmin() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		re := handler.TellMemUsage()
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// Return a random useless fact
	case cmdName == "fact":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		re, err := handler.RandomFact(b.HttpClient)
		if err != nil {
			response.SlashCmdResponseEdit(s, i.Interaction, userObj,
				fmt.Sprintf("An error occured fetching a useless fact: %v", err), true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// Play a sound in the voice channel the requesting user is in
	case cmdName == "play":
		soundName := i.Data.Options[0].StringValue()
		if b.Audio[soundName].Buffer == nil {
			response.SlashCmdResponseMention(s, i.Interaction, userObj,
				fmt.Sprintf("I don't have a registered sound file call %v", soundName))
			return
		}

		var guildObj *discordgo.Guild
		if chanInfo != nil {
			guildObj, err = s.State.Guild(chanInfo.GuildID)
			if err != nil {
				l.Errorf("Could not find guild of channel %q: %v", chanInfo.GuildID, err)
			}
		}
		if guildObj == nil {
			response.SlashCmdResponseMention(s, i.Interaction, userObj,
				"You are not part of a voice channel right now.")
			return
		}

		response.SlashCmdResponseDeferred(s, i.Interaction)
		b.AudioMutex.Lock()
		err := handler.PlaySound(guildObj.VoiceStates, s, *b.Audio[soundName].Buffer, userObj.AuthorId, guildObj.ID)
		if err != nil {
			l.Errorf("An error occured when playing sound: %v", err)
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		b.AudioMutex.Unlock()
		response.SlashCmdDel(s, i.Interaction)
		return

	// SoT: Tell the requester about the currently active daily deed in SoT
	case cmdName == "dailydeed":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsRegistered() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		if !userObj.HasRatCookie() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		em, err := handler.GetDailyDeed(b.HttpClient, userObj, b.Db)
		if err != nil {
			response.SlashCmdResponseEdit(s, i.Interaction, userObj,
				fmt.Sprintf("An error occured fetching the daily deed: %v", err), true)
			return
		}
		response.SlashCmdEmbedDeferred(s, i.Interaction, em)
		return

	// SoT: Reply with the requester's latest achievement in SoT
	case cmdName == "achievement":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsRegistered() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		if !userObj.HasRatCookie() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		em, err := handler.GetSotAchievement(b.HttpClient, userObj)
		if err != nil {
			response.SlashCmdResponseEdit(s, i.Interaction, userObj,
				fmt.Sprintf("An error occured fetching the latest achievement: %v", err), true)
			return
		}
		response.SlashCmdEmbedDeferred(s, i.Interaction, em)
		return

	// SoT: Retrieve the current user balance from the SoT API
	case cmdName == "balance":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsRegistered() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		if !userObj.HasRatCookie() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		re, err := handler.GetSotBalance(b.Db, b.HttpClient, userObj)
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT balance: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// SoT: Retrieve user reputation with a specific faction/company
	case cmdName == "reputation":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsRegistered() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		if !userObj.HasRatCookie() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		repFaction := i.Data.Options[0].StringValue()
		re, err := handler.GetSotReputation(b.HttpClient, userObj, repFaction)
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT reputation: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// SoT: Retrieve user ledger position with a specific faction/company
	case cmdName == "ledger":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsRegistered() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		if !userObj.HasRatCookie() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		ledFaction := i.Data.Options[0].StringValue()
		re, err := handler.GetSotLedger(b.HttpClient, userObj, ledFaction)
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT ledger: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// SoT: Retrieve user ledger position with a specific faction/company
	case cmdName == "stats":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsRegistered() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		if !userObj.HasRatCookie() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		re, err := handler.GetSotStats(b.HttpClient, userObj)
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT user stats: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// SoT: Retrieve user ledger position with a specific faction/company
	case cmdName == "code":
		em, err := handler.GetSotRandomCode()
		if err != nil {
			re := fmt.Sprintf("An error occurred checking your SoT user stats: %v", err.Error())
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdEmbed(s, i.Interaction, em)
		return

	// RareThief/SoT: Fetch and announce the currently active SoT trade routes
	case cmdName == "traderoutes":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		em, err := handler.GetTraderoutes(b.HttpClient, b.Db)
		if err != nil {
			re := fmt.Sprintf("An error occurred retrieving the trade routes: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdEmbedDeferred(s, i.Interaction, em)
		return

	// SoT: Retrieve the users season progress
	case cmdName == "season":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if !userObj.IsRegistered() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		if !userObj.HasRatCookie() {
			response.SlashCmdDel(s, i.Interaction)
			return
		}
		re, err := handler.GetSotSeasonProgress(b.HttpClient, userObj)
		if err != nil {
			re = fmt.Sprintf("An error occurred checking your SoT season progress: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// OWM: Return the weather conditions in a specific location
	case cmdName == "weather":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if b.OwmClient == nil {
			re := "You haven't specified a OpenWeatherMap API key in your config file."
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		weatherLoc := i.Data.Options[0].StringValue()
		re, err := handler.GetCurrentWeather(b.OwmClient, weatherLoc)
		if err != nil {
			re = fmt.Sprintf("An error occurred fetching weather information: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return

	// UD: Return the explanation of a random or specific term from UD
	case cmdName == "urban":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		udTerm := ""
		if len(i.Data.Options) == 1 {
			udTerm = i.Data.Options[0].StringValue()
		}
		em, err := handler.UrbanDict(b.HttpClient, udTerm)
		if err != nil {
			re := fmt.Sprintf("An error occurred fetching term from UD: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdEmbedDeferred(s, i.Interaction, em)
		return

	// TMDB: Return the details of a random or specific movie
	case cmdName == "movie":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if b.TMDb == nil {
			re := "You haven't specified a TMDB API key in your config file."
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		var em *discordgo.MessageEmbed
		var err error
		if len(i.Data.Options) != 1 {
			em, err = handler.TMDbRandMovie(b.TMDb)
		}
		if len(i.Data.Options) == 1 {
			movieTitle := i.Data.Options[0].StringValue()
			em, err = handler.TMDbSearchMovie(b.TMDb, movieTitle)
		}
		if err != nil {
			re := fmt.Sprintf("An error occurred fetching result from TMDB: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdEmbedDeferred(s, i.Interaction, em)
		return

	// TMDB: Return the details of a random or specific TV show
	case cmdName == "tv":
		response.SlashCmdResponseDeferred(s, i.Interaction)
		if b.TMDb == nil {
			re := "You haven't specified a TMDB API key in your config file."
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		var em *discordgo.MessageEmbed
		var err error
		if len(i.Data.Options) != 1 {
			em, err = handler.TMDbRandTvShow(b.TMDb)
		}
		if len(i.Data.Options) == 1 {
			showTitle := i.Data.Options[0].StringValue()
			em, err = handler.TMDbSearchTvShow(b.TMDb, showTitle)
		}
		if err != nil {
			re := fmt.Sprintf("An error occurred fetching result from TMDB: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdEmbedDeferred(s, i.Interaction, em)
		return

	// UserMgmt/SoT: Store the RAT cookie in the Bot DB
	case cmdName == "setrat":
		response.SlashCmdResponseDeferredEphemeral(s, i.Interaction)
		if !userObj.IsRegistered() {
			re := "Sorry, but you are not a registered user. Please contact a admin to get registered."
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		ratCookie := i.Data.Options[0].StringValue()
		re, err := handler.UserSetRatCookie(b.Db, b.Config, userObj, ratCookie)
		if err != nil {
			re = fmt.Sprintf("An error occurred storing your RAT cookie: %v", err)
			response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
			return
		}
		response.SlashCmdResponseEdit(s, i.Interaction, userObj, re, true)
		return
	}
}
