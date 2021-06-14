package bot

import (
	"github.com/bwmarrin/discordgo"
)

// List of slash commands and descriptions
func (b *Bot) SlashCmdList() []*discordgo.ApplicationCommand {

	audioFiles := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for audioFile := range b.Audio {
		audioChoice := &discordgo.ApplicationCommandOptionChoice{
			Name:  audioFile,
			Value: audioFile,
		}
		audioFiles = append(audioFiles, audioChoice)
	}

	slashCmds := []*discordgo.ApplicationCommand{
		// The /help command DMs the requesting user a list of available commands
		{
			Name:        "help",
			Description: "Let SoTBot DM you a list of all available commands",
		},

		// The /version command tells some details about the bot version
		{
			Name:        "version",
			Description: "Let SoTBot tell you some details about itself",
		},

		// With /play you can have the bot join a voice channel and play
		// a pre-defined audio file
		{
			Name:        "play",
			Description: "Let SoTBot join the voice channel you are currently in an play a sound",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "soundfile",
					Description: "Name of the registered sound file",
					Required:    true,
					Choices:     audioFiles,
				},
			},
		},

		// Let the bot tell you how late it is with /time
		{
			Name:        "time",
			Description: "Let SoTBot tell you how late it currently is",
		},

		// SoT: Retrieve the daily deeds from SoT
		{
			Name:        "dailydeed",
			Description: "Let SoTBot show you the currently active daily deed to complete in SoT",
		},
	}

	return slashCmds
}
