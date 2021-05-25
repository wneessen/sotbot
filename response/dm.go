package response

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/user"
)

func DmUser(s *discordgo.Session, u *user.User, msg string, mention bool, noarr bool) {
	l := log.WithFields(log.Fields{
		"action": "bot.DmUser",
	})

	arrText := ""
	if !noarr {
		arrText = fmt.Sprintf("%v!", RandomArrr())
		if mention {
			arrText = fmt.Sprintf("%v %v!", RandomArrr(), u.Mention)
		}
	}

	st, err := s.UserChannelCreate(u.AuthorId)
	if err != nil {
		l.Errorf("Failed to initiate DM channel with user: %v", err)
		return
	}

	_, err = s.ChannelMessageSend(st.ID, fmt.Sprintf("%v %v", arrText, msg))
	if err != nil {
		l.Errorf("Failed to notify user: %v", err)
	}
}
