package nagios

import (
	"fmt"
	"os"
)

// BUG(me): Exit doesn't check for newlines/pipes in msg

// Exit terminates the plugin with the given status and message text, adding
// any performance data that was created.
func Exit(status Status, msg string) {
	fmt.Printf("%v: %s%s\n", status, msg, globalPerfdata)
	os.Exit(int(status))
}

// A Status represents Nagios' interpretation of a return code.
type Status int

const (
	OK Status = iota
	WARNING
	CRITICAL
	UNKNOWN
)

// String returns status as a string.
func (status Status) String() string {
	switch status {
	case OK:
		return "OK"
	case WARNING:
		return "WARNING"
	case CRITICAL:
		return "CRITICAL"
	}
	return "UNKNOWN"
}
