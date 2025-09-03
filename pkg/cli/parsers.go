package cli

import (
	"strconv"
	"strings"
)

func ParseIdsList(arg string) []uint {
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
