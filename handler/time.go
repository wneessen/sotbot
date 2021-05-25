package handler

import (
	"fmt"
	"time"
)

// Let's the bot tell you the current date/time when requested via !time command
func TellTime() string {
	return fmt.Sprintf("The current time is: %v", time.Now().Format("2006-01-02 15:04:05 MST"))
}
