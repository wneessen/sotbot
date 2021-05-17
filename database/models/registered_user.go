package models

type RegisteredUser struct {
	General
	UserId  string
	IsAdmin bool
}
