package handler

import (
	"fmt"
	"time"
)

// Let's the bot provide a help message in the DMs
func Uptime(s time.Time) (string, error) {
	nowUnix := time.Now().Unix()
	timeDiff := nowUnix - s.Unix()
	durTime, err := time.ParseDuration(fmt.Sprintf("%ds", timeDiff))
	if err != nil {
		return "", err
	}
	responseMsg := fmt.Sprintf("I was started: %v, so I am running for %v now",
		s.Format(time.RFC1123), durTime.String())

	return responseMsg, nil
}
