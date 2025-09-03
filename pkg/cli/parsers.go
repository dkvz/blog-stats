package cli

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type CliArgs struct {
	Mode        int
	IgnoredIds  []uint
	StartLength uint
	EndLength   uint
}

func ParseCliArgs() (*CliArgs, error) {
	// TODO: I should use subflags for some of these that have no
	// sense in verify mode.
	mode := flag.String("mode", "", "plot | verify")

	ignoreIds := flag.String(
		"ignore-ids",
		"",
		"comma separated list of article IDs to ignore in computations",
	)
	startLength := flag.Uint("start-length",
		0,
		"only include articles with length higher or equal to this value",
	)
	endLength := flag.Uint("end-length",
		0,
		"only include articles with length lower than this value",
	)

	flag.Parse()

	if (*startLength != 0 && *endLength != 0) && *startLength >= *endLength {
		return nil, errors.New("start-length should be smaller than end-length")
	}
	*mode = strings.TrimSpace(strings.ToLower(*mode))
	ignoredIds := parseIdsList(*ignoreIds)

	iMode := 0

	// My Java past is making me nervous when I don't check for nil
	// values on these pointers but the default val above should
	// make it so they're never nil
	switch *mode {
	case "plot":
		iMode = 0
	case "verify":
		iMode = 1
	default:
		return nil, errors.New("invalid mode")
	}

	return &CliArgs{
		Mode:        iMode,
		IgnoredIds:  ignoredIds,
		StartLength: *startLength,
		EndLength:   *endLength,
	}, nil
}

func parseIdsList(arg string) []uint {
	var ret []uint

	for v := range strings.SplitSeq(arg, ",") {
		conv, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
		if err != nil {
			// Ignore failed parsed values
			continue
		}

		ret = append(ret, uint(conv))
	}

	return ret
}
