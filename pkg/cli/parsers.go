package cli

import (
	"errors"
	"flag"
	"sort"
	"strconv"
	"strings"
)

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
	defaultFactor := flag.Float64("default-factor", 0.0, "default or only factor to use for the most simple linear prediction model WordCount=(factor x length)")
	reg := flag.String("reg", "", "comma separated linear regression params b,a as in y=a+bx (b is the coefficient)")
	// The factor option can be used multiple times
	var factors multiArg
	flag.Var(
		&factors,
		"factor",
		"comma separated values in format factor,range-start,range-end (ranges of article lengths) can be used multiple times - Requires the default-factor option to be present as well",
	)

	flag.Parse()

	if (*startLength != 0 && *endLength != 0) && *startLength >= *endLength {
		return nil, errors.New("start-length should be smaller than end-length")
	}
	*mode = strings.TrimSpace(strings.ToLower(*mode))
	ignoredIds := parseIdsList(*ignoreIds)

	iMode := 0
	var verifyModeArgs *VerifyArgs

	// My Java past is making me nervous when I don't check for nil
	// values on these pointers but the default val above should
	// make it so they're never nil
	switch *mode {
	case "plot":
		iMode = 0
	case "verify":
		// We either need default-factor (+ eventually more factors)
		// or "reg" - We don't accept both of them unless we take
		// one to get precedence.
		var err error
		verifyModeArgs, err = validateVerifyModeArgs(*defaultFactor, factors, *reg)
		if err != nil {
			return nil, err
		}
		iMode = 1
	default:
		return nil, errors.New("invalid mode")
	}

	return &CliArgs{
		Mode:           iMode,
		IgnoredIds:     ignoredIds,
		StartLength:    *startLength,
		EndLength:      *endLength,
		VerifyModeArgs: verifyModeArgs,
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

func validateVerifyModeArgs(
	defaultFactor float64,
	factors multiArg,
	reg string,
) (*VerifyArgs, error) {
	var ret *VerifyArgs

	// If defaultFactor is > 0, use factor mode
	// (default) and check if we have factors and
	// if they're all valid
	if defaultFactor > 0 {
		var parsedFactors []Factor
		for _, f := range factors {
			fact, err := ParseFactor(f)
			if err != nil {
				return nil, err
			}
			// I could also check if any of the ranges overlap
			// but we'll leave them as is and the system will
			// pick the first factor that matches
			parsedFactors = append(parsedFactors, *fact)
		}

		// Sort the factors:
		sort.Slice(parsedFactors, func(i, j int) bool {
			return parsedFactors[i].Start < parsedFactors[j].Start
		})

		ret = &VerifyArgs{
			DefaultFactor: defaultFactor,
			Factors:       parsedFactors,
		}
		return ret, nil
	}

	// Otherwise check if reg can be parsed
	// If not, return missing argument error
	if reg == "" {
		return nil, errors.New("missing 'reg' or 'default-factor' arguments")
	}

	regParms := strings.Split(reg, ",")
	if len(regParms) < 2 {
		return nil, errors.New("format for 'reg' argument is b,a (2 comma separated values) as in y=a+bx")
	}

	b, err := strconv.ParseFloat(strings.TrimSpace(regParms[0]), 64)
	if err != nil {
		return nil, err
	}

	a, err := strconv.ParseFloat(strings.TrimSpace(regParms[1]), 64)
	if err != nil {
		return nil, err
	}

	ret = &VerifyArgs{
		DefaultFactor: b,
		RegA:          a,
		RegMode:       true,
	}

	return ret, nil
}
