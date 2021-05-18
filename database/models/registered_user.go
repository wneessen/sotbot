package models

type RegisteredUser struct {
	General
	UserId string `gorm:"unique"`
}
