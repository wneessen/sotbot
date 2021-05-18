package models

type UserPref struct {
	General
	UserID uint           `gorm:"index:idx_user_key"`
	User   RegisteredUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Key    string         `gorm:"index:idx_user_key"`
	Value  string         `gorm:"size:1024"`
}
