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

		// Report the bots uptime
		{
			Name:        "uptime",
			Description: "Let SoTBot tell you how long it's running so far",
		},

		// Report the bots memory usage
		{
			Name:        "memory",
			Description: "Let SoTBot tell you a bit about it's memory usage",
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

		// SoT: Retrieve the users latest achievement in SoT
		{
			Name:        "achievement",
			Description: "Let SoTBot tell you about your latest achievement in SoT",
		},

		// SoT: Get the users current SoT balance
		{
			Name:        "balance",
			Description: "Let SoTBot tell you your current SoT balance",
		},

		// SoT: Get users faction/company reputation
		{
			Name:        "reputation",
			Description: "Let SoTBot tell you your current SoT reputation with a specific faction/company",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "faction",
					Description: "Name of the faction/company",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Athena's Fortune", Value: "athena"},
						{Name: "Bilge Rats", Value: "bilge"},
						{Name: "Gold Hoarder", Value: "hoarder"},
						{Name: "Hunter's Call", Value: "hunter"},
						{Name: "Merchang Alliance", Value: "merchant"},
						{Name: "Order of Souls", Value: "order"},
						{Name: "Reaper's Bone", Value: "reaper"},
						{Name: "Sea Dogs", Value: "seadog"},
					},
				},
			},
		},

		// SoT: Get users overall ledger position for a specific faction/company
		{
			Name:        "ledger",
			Description: "Let SoTBot tell you your current SoT ledger position with a specific faction/company",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "faction",
					Description: "Name of the faction/company",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Athena's Fortune", Value: "athena"},
						{Name: "Gold Hoarder", Value: "hoarder"},
						{Name: "Merchang Alliance", Value: "merchant"},
						{Name: "Order of Souls", Value: "order"},
						{Name: "Reaper's Bone", Value: "reaper"},
					},
				},
			},
		},

		// SoT: Get the users stats
		{
			Name:        "stats",
			Description: "Let SoTBot tell you some SoT user stats",
		},

		// SoT: Show a random article from the SoT pirate code
		{
			Name:        "code",
			Description: "Let SoTBot tell you a random article from the SoT pirate code",
		},

		// SoT: Show season progress
		{
			Name:        "season",
			Description: "Let SoTBot tell you your current progress in the running SoT season",
		},

		// SoT/Rarethief: Annouce the currently active trading routes
		{
			Name:        "traderoutes",
			Description: "Let SoTBot tell you the currently active trade routes in SoT",
		},
	}

	return slashCmds
}
