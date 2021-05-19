package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/database/models"
	"github.com/wneessen/sotbot/sotapi"
	"gorm.io/gorm"
	"time"
)

func GetBalance(d *gorm.DB, u uint) (models.SotBalance, error) {
	balanceObj := models.SotBalance{}
	d.Where(models.SotBalance{
		UserID: u,
	}).First(&balanceObj)

	if balanceObj.ID <= 0 {
		return balanceObj, fmt.Errorf("No balance information found in database")
	}

	return balanceObj, nil
}

func GetLatestHistoryBalance(d *gorm.DB, u uint) models.SotBalanceHistory {
	historyObj := models.SotBalanceHistory{}
	d.Where(models.SotBalanceHistory{
		UserID: u,
	}).Last(&historyObj)

	return historyObj
}

func GetBalanceDifference(d *gorm.DB, u uint) models.SotBalance {
	returnBalance := models.SotBalance{}
	curBalance, err := GetBalance(d, u)
	if err != nil {
		return returnBalance
	}

	histBalance := GetLatestHistoryBalance(d, u)
	if histBalance.ID <= 0 {
		returnBalance.Gold = 0
		returnBalance.AncientCoins = 0
		returnBalance.Doubloons = 0
		returnBalance.LastUpdated = time.Now().Unix()
		return returnBalance
	}

	returnBalance.Gold = curBalance.Gold - histBalance.Gold
	returnBalance.AncientCoins = curBalance.AncientCoins - histBalance.AncientCoins
	returnBalance.Doubloons = curBalance.Doubloons - histBalance.Doubloons
	returnBalance.LastUpdated = histBalance.LastUpdated

	return returnBalance
}

func UpdateBalance(d *gorm.DB, u uint, b *sotapi.UserBalance) error {
	oldBalance := models.SotBalance{}
	d.Where(models.SotBalance{
		UserID: u,
	}).First(&oldBalance)

	if oldBalance.ID > 0 {
		if oldBalance.Gold != b.Gold ||
			oldBalance.Doubloons != b.Doubloons ||
			oldBalance.AncientCoins != b.AncientCoins {
			log.Debug("Balance has changed. Updating and storing history entry")
			historyBalance := models.SotBalanceHistory{
				UserID:       oldBalance.UserID,
				Gold:         oldBalance.Gold,
				AncientCoins: oldBalance.AncientCoins,
				Doubloons:    oldBalance.Doubloons,
				LastUpdated:  oldBalance.LastUpdated,
			}
			dbTx := d.Create(&historyBalance)
			if dbTx.Error != nil {
				return dbTx.Error
			}
			oldBalance.Gold = b.Gold
			oldBalance.AncientCoins = b.AncientCoins
			oldBalance.Doubloons = b.Doubloons
			oldBalance.LastUpdated = time.Now().Unix()
			dbTx = d.Save(&oldBalance)
			if dbTx.Error != nil {
				return dbTx.Error
			}

			return nil
		}
		log.Debug("Balance didn't change since last check. Skipping update")
		return nil
	}

	balanceEntry := models.SotBalance{
		UserID:       u,
		Gold:         b.Gold,
		AncientCoins: b.AncientCoins,
		Doubloons:    b.Doubloons,
		LastUpdated:  time.Now().Unix(),
	}
	dbTx := d.Create(&balanceEntry)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	return nil
}
