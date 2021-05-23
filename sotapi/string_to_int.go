package sotapi

import (
	"fmt"
	"strconv"
	"strings"
)

type ApiStringInt int64

// UnmarshalJSON function for strings that are obviously ints
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
