package handler

import (
	"fmt"
	"time"
)

// Let's the bot provide a help message in the DMs
func Uptime(s time.Time) (string, bool, error) {
	nowUnix := time.Now().Unix()
	timeDiff := nowUnix - s.Unix()
	durTime, err := time.ParseDuration(fmt.Sprintf("%ds", timeDiff))
	if err != nil {
		return "", false, err
	}
	responseMsg := fmt.Sprintf("I was started: %v, so I am running for %v now",
		s.Format("2006-01-02 15:04:05"), durTime.String())

	return responseMsg, true, nil
}
