// Package nagios provides primitives for writing Nagios plugins.
package nagios

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// A Range is for checking whether a given value is contained within it.
// Nagios plugins conventionally expect textual representations of ranges to be
// of the format [@]start:end, where:
//	1. start <= end
//	2. "start:" is optional if start=0
//	3. if "start:" is present and end isn't specified, then end is infinity
//	4. "~" represents negative infinity
//	5. if "@" is present, then range is the complement of start:end
type Range struct {
	complement bool
	start      float64
	end        float64
}

// NewRange allocates and sets a new Range based on s, returning its address.
func NewRange(s string) (*Range, error) {
	r := new(Range)
	err := r.Set(s)
	return r, err
}

// InRange returns true if v is in r.
func (r *Range) InRange(v float64) bool {
	if r.complement {
		return v < r.start || r.end < v
	}
	return r.start <= v && v <= r.end
}

// String returns a textual representation of r.
func (r *Range) String() string {
	var s string

	// return empty string if r is nil
	if r == nil {
		return s
	}

	// begin the string with a "@" if r.complement is set
	if r.complement {
		s = "@"
	}

	// print the start value unless it's zero (although ranges like @20 are
	// confusing, make the zero explicit in that case)
	if math.IsInf(r.start, -1) {
		s += "~:"
	} else if r.start != 0 {
		s += strconv.FormatFloat(r.start, 'f', -1, 64) + ":"
	} else if r.complement {
		s += "0:"
	}

	// print the end value unless it's +Inf, making sure we don't end up
	// with an empty string
	if !math.IsInf(r.end, 1) {
		s += strconv.FormatFloat(r.end, 'f', -1, 64)
	} else if len(s) == 0 {
		s += "0:"
	}

	return s
}

// Set sets r to the Range described by s.
func (r *Range) Set(s string) error {
	var err error
	var lower, upper string

	// make sure the range isn't empty
	if s == "" || s == "@" {
		return errors.New("empty range")
	}

	// set r.complement and remove "@" from the string if it's present
	if s[0] == '@' {
		r.complement = true
		s = s[1:]
	}

	// split the input and assign each part to variables
	if endpoints := strings.SplitN(s, ":", 2); len(endpoints) == 1 {
		lower, upper = "0", endpoints[0]
	} else {
		lower, upper = endpoints[0], endpoints[1]
	}

	// change "~" to "-Inf" so Go will interpret it as negative infinity
	if lower == "~" {
		lower = "-Inf"
	}

	// parse the start value
	r.start, err = strconv.ParseFloat(lower, 64)
	if err != nil {
		return err
	}

	// change an empty upper string to "Inf" to be interpreted as infinity
	if upper == "" {
		upper = "Inf"
	}

	// parse the end value
	r.end, err = strconv.ParseFloat(upper, 64)
	if err != nil {
		return err
	}

	// check for invalid ranges
	if r.start+r.end == math.NaN() {
		return errors.New("NaN not allowed in ranges")
	} else if r.start > r.end {
		return errors.New("start greater than end")
	}

	return nil
}
