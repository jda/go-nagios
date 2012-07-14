package nagios

import (
	"errors"
	"fmt"
	"strings"
)

// A PerfData represents one data point to be graphed by third-party programs.
// Unit may be one of:
//	"" - none
//	s, us, ms - seconds, microseconds, milliseconds
//	% - percent
//	B, kB, MB, GB, TB - bytes, kilobytes, megabytes, gigabytes, terabytes
//	c - continuous counter
type PerfData struct {
	Label      string
	Value      float64
	Unit       string
	Warn, Crit *Range
	Min, Max   *float64
}

// perfdata tracks everything added by AddPerfData to be printed with Exit
var perfdata string

// BUG(me): AddPerfData doesn't check for newlines in label

// AddPerfData adds p to the global perfdata string to be printed with Exit
func AddPerfData(p PerfData) error {
	// check for valid unit of measure
	switch strings.ToLower(p.Unit) {
	case "", "s", "us", "ms", "%", "b", "kb", "mb", "gb", "tb", "c":
	default:
		return errors.New("invalid unit of measure")
	}

	// append a pipe or space character, depending on whether or not this
	// is the first perfdata added
	if len(perfdata) == 0 {
		perfdata = "|"
	} else {
		perfdata += " "
	}

	// add label, value, unit of measure, and warning/critical ranges
	perfdata += fmt.Sprintf("%s=%v%s;%v;%v", escapeLabel(p.Label),
		p.Value, p.Unit, p.Warn, p.Crit)

	// add min and max values if they are non-nil
	// unnecessary if measured as a percentage
	if p.Unit != "%" {
		for _, v := range []*float64{p.Min, p.Max} {
			perfdata += ";"
			if v != nil {
				perfdata += fmt.Sprintf("%v", *v)
			}
		}
	}

	// remove trailing semicolons
	perfdata = strings.TrimRight(perfdata, ";")

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
