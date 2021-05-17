package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(d, ll string) (*gorm.DB, error) {
	l := log.WithFields(log.Fields{
		"action": "database.ConnectDB",
	})

	dbLogLevel := logger.LogLevel(1)
	if ll == "debug" {
		dbLogLevel = logger.LogLevel(2)
	}
	db, err := gorm.Open(sqlite.Open(d), &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.RegisteredUser{}); err != nil {
		l.Errorf("Database automigration failed: %v", err)
	}

	return db, nil
}
