package nagios

import (
	"errors"
	"fmt"
	"strings"
)

// Unit may be one of:
//	"" - none
//	s, us, ms - seconds, microseconds, milliseconds
//	% - percent
//	B, kB, MB, GB, TB - bytes, kilobytes, megabytes, gigabytes, terabytes
//	c - continuous counter

// perfdata tracks everything added by Perfdata to be printed with Exit
var globalPerfdata string

// Perfdata adds p to the global perfdata string to be printed with Exit
func Perfdata(label string, v float64, unit string, warn, crit *Range,
	extrema ...float64) error {

	// check that extrema has no more than 2 values, and set min and max
	// values from it
	if len(extrema) > 2 {
		return errors.New("too many arguments")
	}

	// check for valid unit of measure
	switch strings.ToLower(unit) {
	case "", "s", "us", "ms", "%", "b", "kb", "mb", "gb", "tb", "c":
	default:
		return errors.New("invalid unit of measure " + unit)
	}

	// check for (some) invalid characters in label
	if strings.ContainsRune(label, '\n') {
		return errors.New("label contains invalid characters")
	}

	// append a pipe or space character, depending on whether or not this
	// is the first perfdata added
	if len(globalPerfdata) == 0 {
		globalPerfdata = "|"
	} else {
		globalPerfdata += " "
	}

	// add label, value, unit of measure, and warning/critical ranges
	globalPerfdata += fmt.Sprintf("%s=%v%s;%v;%v", escapeLabel(label), v,
		unit, warn, crit)

	// add min and max values if they are non-nil
	// unnecessary if measured as a percentage
	if unit != "%" {
		for _, v := range extrema {
			globalPerfdata += fmt.Sprintf(";%v", v)
		}
	}

	// remove trailing semicolons
	globalPerfdata = strings.TrimRight(globalPerfdata, ";")

	return nil
}

// escapeLabel makes label safe for use in a perfdata string by wrapping it in
// single quotes if it contains a space, single quote, or equals character,
// doubling any single quotes
func escapeLabel(label string) string {
	if strings.ContainsAny(label, " ='") {
		return "'" + strings.Replace(label, "'", "''", -1) + "'"
	}
	return label
}
