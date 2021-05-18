package models

type SotBalance struct {
	General
	UserID       uint           `gorm:"unique"`
	User         RegisteredUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Gold         int
	Doubloons    int
	AncientCoins int
	LastUpdated  int64
}

type SotBalanceHistory struct {
	General
	UserID       uint
	User         RegisteredUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Gold         int
	Doubloons    int
	AncientCoins int
	LastUpdated  int64
}
