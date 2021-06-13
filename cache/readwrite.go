package cache

import (
	"fmt"
	"github.com/wneessen/sotbot/database"
	"gorm.io/gorm"
)

func Store(k string, o interface{}, d *gorm.DB) error {
	objString, err := SerializeObj(o)
	if err != nil {
		return err
	}
	if err := database.StoreBotCache(d, k, objString); err != nil {
		return err
	}

	return nil
}

func Read(k string, o interface{}, d *gorm.DB) error {
	objString := database.ReadBotCache(d, k)
	if objString == "" {
		return fmt.Errorf("Object not found in bot cache")
	}
	if err := DeserializeObj(objString, o); err != nil {
		return err
	}

	return nil
}
