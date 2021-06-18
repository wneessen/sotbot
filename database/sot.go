package database

import (
	"gorm.io/gorm"
	"strconv"
)

func GetFailedRatCookieTries(d *gorm.DB, u uint) (int64, error) {
	failedRatTries := UserGetPrefString(d, u, "failed_rat_tries")
	if failedRatTries == "" {
		return 0, nil
	}

	failedRatTriesNum, err := strconv.ParseInt(failedRatTries, 10, 32)
	if err != nil {
		return 0, err
	}

	return failedRatTriesNum, nil
}

func IncreaseFailedRatCookieTries(d *gorm.DB, u uint) (int64, error) {
	failedCounter, err := GetFailedRatCookieTries(d, u)
	if err != nil {
		return 0, err
	}
	failedCounter++
	failedCounterString := strconv.FormatInt(failedCounter, 10)

	if err := UserSetPref(d, u, "failed_rat_tries", failedCounterString); err != nil {
		return 0, err
	}

	return failedCounter, nil
}
