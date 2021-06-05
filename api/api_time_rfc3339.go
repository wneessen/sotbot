package api

import (
	"fmt"
	"strings"
	"time"
)

type ApiTimeRFC3339 time.Time

// UnmarshalJSON function for strings that are in RFC3339 format
func (t *ApiTimeRFC3339) UnmarshalJSON(s []byte) error {
	dateString := string(s)
	dateString = strings.ReplaceAll(dateString, `"`, "")
	if dateString == "null" {
		return nil
	}
	dateParse, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return fmt.Errorf("failed to parse string as RFC3339 time string: %v", err)
	}

	*(*time.Time)(t) = dateParse
	return nil
}

func (t ApiTimeRFC3339) Time() time.Time {
	return time.Time(t)
}
