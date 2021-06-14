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

	userObj, err := user.NewUser(b.Db, b.Config, i.Member.User.ID)
	if err != nil {
		l.Errorf("Could not create user object: %v", err)
		return
	}
	userObj.AuthorName = i.Member.User.Username
	userObj.Mention = i.Member.User.Mention()

	if chanInfo.Type != discordgo.ChannelTypeDM {
		userPerm, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			l.Errorf("Failed to look up user-channel permission: %v", err)
		}
		userObj.ChanPermission = userPerm
	}

	switch {
	case cmdName == "help":
		helpTexts := handler.Help()
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf("%v, please check your DMs", userObj.Mention),
			},
		}); err != nil {
			l.Errorf("Failed to respond to interaction")
		}
		for _, textPart := range helpTexts {
			response.DmUser(s, userObj, "`"+textPart+"`", false, true)
		}
		return

	case cmdName == "version":
		versionString := handler.TellVersion()
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: versionString,
			},
		}); err != nil {
			l.Errorf("Failed to execute version string slash command: %v", err)
		}
		return

	case cmdName == "play":
		soundName := i.Data.Options[0].StringValue()
		if b.Audio[soundName].Buffer == nil {
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf("I don't have a soundfile called %q", soundName),
				},
			}); err != nil {
				l.Errorf("Failed respond to user's slash command request: %v", err)
			}
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
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: "You are not part of a voice channel right now.",
				},
			}); err != nil {
				l.Errorf("Failed respond to user's slash command request: %v", err)
			}
			return
		}

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "",
			},
		}); err != nil {
			l.Errorf("Failed respond to user's slash command request: %v", err)
			return
		}

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
	}
}
