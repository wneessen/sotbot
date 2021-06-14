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
		response.SlashCmdResponse(s, i.Interaction, userObj, "Please check your DMs.", true)
		for _, textPart := range helpTexts {
			response.DmUser(s, userObj, "`"+textPart+"`", false, true)
		}
		return

	// Reply with the version string of the bot
	case cmdName == "version":
		versionString := handler.TellVersion()
		response.SlashCmdResponse(s, i.Interaction, userObj, versionString, true)
		return

	// Tell the current time
	case cmdName == "time":
		timeString := handler.TellTime()
		response.SlashCmdResponse(s, i.Interaction, userObj, timeString, true)
		return

	// Play a sound in the voice channel the requesting user is in
	case cmdName == "play":
		soundName := i.Data.Options[0].StringValue()
		if b.Audio[soundName].Buffer == nil {
			response.SlashCmdResponse(s, i.Interaction, userObj,
				fmt.Sprintf("I don't have a registered sound file call %v", soundName), true)
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
			response.SlashCmdResponse(s, i.Interaction, userObj, "You are not part of a voice channel right now.",
				true)
			return
		}

		response.SlashCmdResponseDeferred(s, i.Interaction)
		b.AudioMutex.Lock()
		err := handler.PlaySound(guildObj.VoiceStates, s, *b.Audio[soundName].Buffer, userObj.AuthorId, guildObj.ID)
		if err != nil {
			l.Errorf("An error occured when playing sound: %v", err)
			if err := s.InteractionResponseDelete(s.State.User.ID, i.Interaction); err != nil {
				l.Errorf("Failed to delete interaction response: %v", err)
				return
			}
			return
		}
		b.AudioMutex.Unlock()
		if err := s.InteractionResponseDelete(s.State.User.ID, i.Interaction); err != nil {
			l.Errorf("Failed to delete interaction response: %v", err)
			return
		}
		return

	// SoT: Tell the requester about the currently active daily deed in SoT
	case cmdName == "dailydeed":
		if !userObj.IsRegistered() {
			return
		}
		if !userObj.HasRatCookie() {
			return
		}
		response.SlashCmdResponseDeferred(s, i.Interaction)
		em, err := handler.GetDailyDeed(b.HttpClient, userObj, b.Db)
		if err != nil {
			if err := s.InteractionResponseDelete(s.State.User.ID, i.Interaction); err != nil {
				l.Errorf("Failed to delete interaction response: %v", err)
				return
			}
			response.SlashCmdResponse(s, i.Interaction, userObj,
				fmt.Sprintf("An error occured fetching the daily deed: %v", err), true)
			return
		}

		response.SlashCmdEmbedDeferred(s, i, em)
		return

	// SoT: Reply with the requester's latest achievement in SoT
	case cmdName == "achievement":
		if !userObj.IsRegistered() {
			return
		}
		if !userObj.HasRatCookie() {
			return
		}
		response.SlashCmdResponseDeferred(s, i.Interaction)
		em, err := handler.GetSotAchievement(b.HttpClient, userObj)
		if err != nil {
			if err := s.InteractionResponseDelete(s.State.User.ID, i.Interaction); err != nil {
				l.Errorf("Failed to delete interaction response: %v", err)
				return
			}
			response.SlashCmdResponse(s, i.Interaction, userObj,
				fmt.Sprintf("An error occured fetching the latest achievement: %v", err), true)
			return
		}

		response.SlashCmdEmbedDeferred(s, i, em)
		return
	}
}
