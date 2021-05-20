package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/random"
)

type SotCodeArticle struct {
	Number      int
	Title       string
	Description string
}

// Let's the bot tell you the current date/time when requested via !time command
func (b *Bot) RandSotCode(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := log.WithFields(log.Fields{
		"action": "handler.RandSotCode",
	})

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.Content == "!code" {
		l.Debugf("Received '!code' request from user %v", m.Author.Username)
		sotCode := []SotCodeArticle{
			{Number: 1, Title: "The Sea Calls To Us All",
				Description: "Everyone is welcome on the Sea of Thieves regardless of age, gender, race, sexuality, " +
					"nationality or creed."},
			{Number: 2, Title: "The Sea Unites Us as One Community",
				Description: "Outside the heat of battle or piracy on the high seas, all crews shall bond together as " +
					"a community of like-minded souls."},
			{Number: 3, Title: "Disputes Are Settled upon the Waves",
				Description: "None shall quarrel or overly dissent against another crew, but let every engagement be " +
					"settled by sword, pistol and good seamanship."},
			{Number: 4, Title: "All Crewmates Are Equal",
				Description: "Let each crewmate be respected as equal and free to follow their own bearing, speak " +
					"openly and vote in affairs of the voyage."},
			{Number: 5, Title: "The Crew Bond Is Sacred",
				Description: "Those who betray their crew and ship through griefing or trolling shall be sent to " +
					"the brig."},
			{Number: 6, Title: "Respect New Pirates and Their Voyage Ahead",
				Description: "May the old legends help to forge new ones: treat new pirates with respect and share " +
					"your knowledge."},
			{Number: 7, Title: "Those Who Cheat Shall Be Punished",
				Description: "Pirates who show bad form and cheat their crew or others shall surely face bitter " +
					"hardships and punishments."},
		}
		sotCodeNum := len(sotCode)
		randNum, err := random.Number(sotCodeNum)
		if err != nil {
			l.Errorf("Failed to generate random number: %v", err)
			return
		}

		messageEmbed := discordgo.MessageEmbed{
			Title: fmt.Sprintf("From the SoT pirates code. Article %d: %v", sotCode[randNum].Number,
				sotCode[randNum].Title),
			Description: sotCode[randNum].Description,
			Type:        discordgo.EmbedTypeArticle,
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &messageEmbed)
		if err != nil {
			l.Errorf("Failed to send embeded message: %v", err)
			return
		}
	}
}
