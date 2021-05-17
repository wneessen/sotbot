package models

type UserPref struct {
	General
	UserID uint
	User   RegisteredUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Key    string
	Value  string `gorm:"size:1024"`
}
