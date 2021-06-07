package api

import (
	"fmt"
	"strings"
	"time"
)

// ApiTimeRFC3339 is an alias type for time.Time
type ApiTimeRFC3339 time.Time

// UnmarshalJSON converts a API RFC3339 formated date strings into a
// time.Time object
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

// Time returns the time.Time object of a ApiTimeRFC3339 type
func (t ApiTimeRFC3339) Time() time.Time {
	return time.Time(t)
}
