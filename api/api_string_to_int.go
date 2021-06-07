package api

import (
	"fmt"
	"strconv"
	"strings"
)

// ApiStringInt is a type alias for a int64. It is meant to be used for
// example in API responses where an obvious Int is returned as string
// in the JSON
type ApiStringInt int64

// UnmarshalJSON converts the ApiStringInt string into an int64
func (s *ApiStringInt) UnmarshalJSON(is []byte) error {
	intString := string(is)
	intString = strings.ReplaceAll(intString, `"`, ``)
	realInt, err := strconv.ParseInt(intString, 10, 64)
	if err != nil {
		return fmt.Errorf("string to int conversion failed: %v", err)
	}
	*(*int64)(s) = realInt

	return nil
}
