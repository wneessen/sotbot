package models

type BotCache struct {
	General
	Key   string `gorm:"index:idx_cache_key"`
	Value string `gorm:"size:8196"`
}
