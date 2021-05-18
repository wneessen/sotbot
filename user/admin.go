package user

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func IsAdmin(s *discordgo.Session, uid, cid string) bool {
	l := log.WithFields(log.Fields{
		"action": "user.IsAdmin",
	})

	userPerm, err := s.UserChannelPermissions(uid, cid)
	if err != nil {
		l.Errorf("Failed to look up user-channel permission: %v", err)
		return false
	}

	return userPerm&discordgo.PermissionAdministrator != 0
}
