package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Factor struct {
	Value float64
	Start uint
	End   uint
}

type VerifyArgs struct {
	// Either the default factor to use or the
	// coefficient for the line equation ("b" as
	// in y=a+bx)
	DefaultFactor float64
	// Means we're using linear regression when true
	// Also implies the default mode is using factors
	RegMode bool
	Factors []Factor
	RegA    float64
	// Outputs an HTML table in verify mode when true
	TableMode bool
}

type CliArgs struct {
	Mode           int
	IgnoredIds     []uint
	StartLength    uint
	EndLength      uint
	VerifyModeArgs *VerifyArgs
}

func ParseFactor(f string) (*Factor, error) {
	vals := strings.Split(f, ",")
	if len(vals) != 3 {
		return nil, errors.New("missing parameters in factor, 3 values expected")
	}

	// First item is supposed to be a float
	fact, err := strconv.ParseFloat(strings.TrimSpace(vals[0]), 64)
	if err != nil {
		return nil, err
	}

	// Next we got two uint
	start, err := strconv.ParseUint(strings.TrimSpace(vals[1]), 10, 32)
	if err != nil {
		return nil, err
	}
	end, err := strconv.ParseUint(strings.TrimSpace(vals[2]), 10, 32)
	if err != nil {
		return nil, err
	}

	if start >= end {
		return nil, errors.New("the start length cannot be equal or higher than the end length")
	}

	return &Factor{
		Value: fact,
		Start: uint(start),
		End:   uint(end),
	}, nil
}

// Custom flag type backed by a simple slice of strings
// We'll be able to provide that arg multiple time
// Feel like this should exist in the lib
type multiArg []string

func (m *multiArg) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *multiArg) Set(v string) error {
	*m = append(*m, v)
	return nil
}
