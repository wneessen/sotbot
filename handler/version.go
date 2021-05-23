package handler

import (
	"fmt"
	"github.com/wneessen/sotbot/version"
)

// Let the bot tell us it's version information
func TellVersion() (string, bool) {
	return fmt.Sprintf("I am SoTBot Version v%v (OS: %v, Arch: %v). I was built at at: %v by: %v)",
		version.Version, version.BuildOs, version.BuildArch, version.BuildDate, version.BuildUser), true
}
