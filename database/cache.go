package database

import (
	"github.com/wneessen/sotbot/database/models"
	"gorm.io/gorm"
)

func StoreBotCache(d *gorm.DB, k, v string) error {
	botCache := models.BotCache{}
	d.Where(models.BotCache{
		Key: k,
	}).First(&botCache)

	if botCache.ID <= 0 {
		dbTx := d.Create(&models.BotCache{
			Key:   k,
			Value: v,
		})
		if dbTx.Error != nil {
			return dbTx.Error
		}
		return nil
	}
	botCache.Value = v
	dbTx := d.Save(&botCache)
	if dbTx.Error != nil {
		return dbTx.Error
	}
	return nil
}

func ReadBotCache(d *gorm.DB, k string) string {
	botCache := models.BotCache{}
	d.Where(models.BotCache{
		Key: k,
	}).First(&botCache)

	return botCache.Value
}
