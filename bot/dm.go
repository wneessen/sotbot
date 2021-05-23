package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/user"
)

func DmUser(s *discordgo.Session, u *user.User, msg string, mention bool) {
	l := log.WithFields(log.Fields{
		"action": "bot.DmUser",
	})

	st, err := s.UserChannelCreate(u.AuthorId)
	if err != nil {
		l.Errorf("Failed to initiate DM channel with user: %v", err)
		return
	}

	if mention {
		_, err = s.ChannelMessageSend(st.ID, RandomArrr()+" "+u.Mention+"! "+msg)
		if err != nil {
			l.Errorf("Failed to notify user: %v", err)
		}
		return
	}
	_, err = s.ChannelMessageSend(st.ID, RandomArrr()+"! "+msg)
	if err != nil {
		l.Errorf("Failed to notify user: %v", err)
	}
}
